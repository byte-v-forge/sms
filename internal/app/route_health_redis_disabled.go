package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *RedisRouteHealthStore) DisabledRouteKeys(ctx context.Context, routes []core.Route) (map[string]struct{}, error) {
	if !s.available() {
		return map[string]struct{}{}, nil
	}
	routeKeys, redisKeys := disabledRouteRedisKeys(routes)
	if len(redisKeys) == 0 {
		return map[string]struct{}{}, nil
	}
	values, err := s.client.MGet(ctx, redisKeys...).Result()
	if err != nil {
		return nil, err
	}
	return disabledRouteKeySet(routeKeys, values), nil
}

func disabledRouteKeySet(routeKeys []string, values []any) map[string]struct{} {
	out := make(map[string]struct{}, len(values))
	for index, raw := range values {
		if raw == nil {
			continue
		}
		out[routeKeys[index]] = struct{}{}
	}
	return out
}
