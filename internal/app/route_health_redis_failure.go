package app

import (
	"context"
	"strconv"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *RedisRouteHealthStore) RecordAcquireFailure(ctx context.Context, route core.Route, err *core.Error) error {
	if !s.available() {
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
