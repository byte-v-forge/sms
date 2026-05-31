package app

import (
	"context"
	"strconv"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *RedisRouteHealthStore) DisabledRouteKeys(ctx context.Context, routes []core.Route) (map[string]struct{}, error) {
	if s == nil || s.client == nil {
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
	out := make(map[string]struct{}, len(values))
	for index, raw := range values {
		if raw == nil {
			continue
		}
		out[routeKeys[index]] = struct{}{}
	}
	return out, nil
}

func (s *RedisRouteHealthStore) RecordAcquireFailure(ctx context.Context, route core.Route, err *core.Error) error {
	if s == nil || s.client == nil {
		return nil
	}
	if !routeFailureCountsTowardDisable(err) {
		return nil
	}
	routeKey := routeHealthKey(route)
	if routeKey == "" {
		return nil
	}
	policy := routeFailurePolicyWithDefaults(route.FailurePolicy)
	return routeFailureScript.Run(
		ctx,
		s.client,
		[]string{routeHealthRedisKey("failure", routeKey), routeHealthRedisKey("disabled", routeKey)},
		strconv.Itoa(policy.FailureThreshold),
		strconv.Itoa(seconds(policy.FailureWindow)),
		strconv.Itoa(seconds(policy.DisableTTL)),
	).Err()
}

func (s *RedisRouteHealthStore) RecordAcquireSuccess(ctx context.Context, route core.Route) error {
	if s == nil || s.client == nil {
		return nil
	}
	routeKey := routeHealthKey(route)
	if routeKey == "" {
		return nil
	}
	return s.client.Del(ctx, routeHealthRedisKey("failure", routeKey), routeHealthRedisKey("disabled", routeKey)).Err()
}
