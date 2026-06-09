package main

import (
	"context"

	"github.com/byte-v-forge/sms/internal/app"
	providerspi "github.com/byte-v-forge/sms/internal/providers/spi"
)

type runtimeStores struct {
	configs     app.ProviderConfigStore
	orders      app.OrderStore
	orderList   app.OrderListStore
	orderOutbox *app.PostgresOrderStore
	codeSecrets *app.SMSCodeSecretStore
	routeHealth app.RouteHealthStore
	close       func()
}

func newRuntimeStores(ctx context.Context, cfg config, providers *providerspi.Registry) (*runtimeStores, error) {
	closers := []func(){}
	store := &runtimeStores{close: func() { closeAll(closers) }}
	clock := app.SystemClock{}

	postgresPool, closePostgres, err := newOptionalPostgresPool(ctx, cfg)
	if err != nil {
		return nil, err
	}
	closers = append(closers, closePostgres)

	configs, err := newProviderConfigStore(ctx, postgresPool, providers, clock)
	if err != nil {
		closeAll(closers)
		return nil, err
	}
	store.configs = configs

	orderHistory, orderOutbox, err := newOrderHistoryStore(ctx, postgresPool)
	if err != nil {
		closeAll(closers)
		return nil, err
	}
	store.orderList = orderHistory
	store.orderOutbox = orderOutbox
	store.orders = orderHistory

	redisClient, closeRedis, err := newOptionalRedisClient(ctx, cfg)
	if err != nil {
		closeAll(closers)
		return nil, err
	}
	closers = append(closers, closeRedis)
	configureRuntimeCacheStores(store, redisClient, orderHistory, clock)
	return store, nil
}

func closeAll(closers []func()) {
	for i := len(closers) - 1; i >= 0; i-- {
		if closers[i] != nil {
			closers[i]()
		}
	}
}
