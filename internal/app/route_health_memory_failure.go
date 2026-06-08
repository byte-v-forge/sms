package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

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
