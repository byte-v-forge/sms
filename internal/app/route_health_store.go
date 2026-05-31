package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/redis/go-redis/v9"
)

type RouteHealthStore interface {
	DisabledRouteKeys(ctx context.Context, routes []core.Route) (map[string]struct{}, error)
	RecordAcquireFailure(ctx context.Context, route core.Route, err *core.Error) error
	RecordAcquireSuccess(ctx context.Context, route core.Route) error
}

type noopRouteHealthStore struct{}

func (noopRouteHealthStore) DisabledRouteKeys(context.Context, []core.Route) (map[string]struct{}, error) {
	return map[string]struct{}{}, nil
}

func (noopRouteHealthStore) RecordAcquireFailure(context.Context, core.Route, *core.Error) error {
	return nil
}

func (noopRouteHealthStore) RecordAcquireSuccess(context.Context, core.Route) error {
	return nil
}

type RedisRouteHealthStore struct {
	client redis.Cmdable
}

func NewRedisRouteHealthStore(client redis.Cmdable) *RedisRouteHealthStore {
	return &RedisRouteHealthStore{client: client}
}
