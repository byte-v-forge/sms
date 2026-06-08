package app

import (
	"context"
	"sync"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

type MemoryRouteHealthStore struct {
	mu       sync.Mutex
	failures map[string]memoryRouteFailure
	disabled map[string]time.Time
	clock    core.Clock
}

type memoryRouteFailure struct {
	count     int
	expiresAt time.Time
}

func NewMemoryRouteHealthStore(clock core.Clock) *MemoryRouteHealthStore {
	if clock == nil {
		clock = SystemClock{}
	}
	return &MemoryRouteHealthStore{
		failures: map[string]memoryRouteFailure{},
		disabled: map[string]time.Time{},
		clock:    clock,
	}
}

func (s *MemoryRouteHealthStore) DisabledRouteKeys(ctx context.Context, routes []core.Route) (map[string]struct{}, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	now := s.clock.Now()
	out := map[string]struct{}{}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, route := range routes {
		routeKey := routeHealthKey(route)
		if routeKey == "" {
			continue
		}
		expiresAt, ok := s.disabled[routeKey]
		if !ok {
			continue
		}
		if routeHealthExpired(now, expiresAt) {
			delete(s.disabled, routeKey)
			continue
		}
		out[routeKey] = struct{}{}
	}
	return out, nil
}

func (s *MemoryRouteHealthStore) RecordAcquireFailure(ctx context.Context, route core.Route, err *core.Error) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !routeFailureCountsTowardDisable(err) {
		return nil
	}
	routeKey := routeHealthKey(route)
	if routeKey == "" {
		return nil
	}
	policy := routeFailurePolicyWithDefaults(route.FailurePolicy)
	now := s.clock.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	failure := s.failures[routeKey]
	if failure.expiresAt.IsZero() || routeHealthExpired(now, failure.expiresAt) {
		failure = memoryRouteFailure{expiresAt: now.Add(policy.FailureWindow)}
	}
	failure.count++
	if failure.count >= policy.FailureThreshold {
		s.disabled[routeKey] = now.Add(policy.DisableTTL)
		delete(s.failures, routeKey)
		return nil
	}
	s.failures[routeKey] = failure
	return nil
}

func (s *MemoryRouteHealthStore) RecordAcquireSuccess(ctx context.Context, route core.Route) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	routeKey := routeHealthKey(route)
	if routeKey == "" {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.failures, routeKey)
	delete(s.disabled, routeKey)
	return nil
}

func routeHealthExpired(now time.Time, expiresAt time.Time) bool {
	return !expiresAt.IsZero() && !now.Before(expiresAt)
}
