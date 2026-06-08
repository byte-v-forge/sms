package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

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
