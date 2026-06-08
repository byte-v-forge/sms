package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

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
