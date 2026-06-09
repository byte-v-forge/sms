package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *RedisRouteHealthStore) RecordAcquireSuccess(ctx context.Context, route core.Route) error {
	if !s.available() {
		return nil
	}
	routeKey := routeHealthKey(route)
	if routeKey == "" {
		return nil
	}
	return s.client.Del(ctx, routeHealthRedisKey("failure", routeKey), routeHealthRedisKey("disabled", routeKey)).Err()
}
