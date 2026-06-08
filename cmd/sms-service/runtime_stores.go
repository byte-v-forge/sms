package main

import (
	"context"
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/app"
	"github.com/byte-v-forge/sms/internal/platform/redisx"
	providerspi "github.com/byte-v-forge/sms/internal/providers/spi"
	"github.com/redis/go-redis/v9"
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

type orderHistoryStore interface {
	app.OrderStore
	app.OrderListStore
}

func newRuntimeStores(ctx context.Context, cfg config, providers *providerspi.Registry) (*runtimeStores, error) {
	closers := []func(){}
	store := &runtimeStores{close: func() { closeAll(closers) }}
	clock := app.SystemClock{}

	configs, closeConfigs, err := newProviderConfigStore(ctx, cfg, providers, clock)
	if err != nil {
		return nil, err
	}
	closers = append(closers, closeConfigs)
	store.configs = configs

	orderHistory, orderOutbox, closeOrders, err := newOrderHistoryStore(ctx, cfg)
	if err != nil {
		closeAll(closers)
		return nil, err
	}
	closers = append(closers, closeOrders)
	store.orderList = orderHistory
	store.orderOutbox = orderOutbox
	store.orders = orderHistory

	redisClient, closeRedis, err := newOptionalRedisClient(ctx, cfg)
	if err != nil {
		closeAll(closers)
		return nil, err
	}
	closers = append(closers, closeRedis)
	if redisClient != nil {
		activeStore := app.NewRedisOrderStore(redisx.NewStringStore(redisClient, "sms:order", 30*time.Minute), clock)
		store.orders = app.NewCompositeOrderStore(activeStore, orderHistory)
		store.codeSecrets = app.NewSMSCodeSecretStore(redisx.NewStringStore(redisClient, "sms:code-secret", 30*time.Minute), clock)
		store.routeHealth = app.NewRedisRouteHealthStore(redisClient)
	} else {
		store.codeSecrets = app.NewSMSCodeSecretStore(app.NewMemoryTTLStringStore("sms:code-secret", 30*time.Minute, clock), clock)
		store.routeHealth = app.NewMemoryRouteHealthStore(clock)
	}
	return store, nil
}

func newProviderConfigStore(ctx context.Context, cfg config, providers *providerspi.Registry, clock app.SystemClock) (app.ProviderConfigStore, func(), error) {
	if strings.TrimSpace(cfg.PGDSN) == "" {
		return app.NewMemoryProviderConfigStore(providers, clock), noopClose, nil
	}
	store, err := app.NewPostgresProviderConfigStore(ctx, cfg.PGDSN, providers)
	if err != nil {
		return nil, nil, err
	}
	return store, store.Close, nil
}

func newOrderHistoryStore(ctx context.Context, cfg config) (orderHistoryStore, *app.PostgresOrderStore, func(), error) {
	if strings.TrimSpace(cfg.PGDSN) == "" {
		store := app.NewMemoryOrderStore()
		return store, nil, noopClose, nil
	}
	store, err := app.NewPostgresOrderStore(ctx, cfg.PGDSN)
	if err != nil {
		return nil, nil, nil, err
	}
	return store, store, store.Close, nil
}

func newOptionalRedisClient(ctx context.Context, cfg config) (*redis.Client, func(), error) {
	rawURL := strings.TrimSpace(cfg.RedisURL)
	if rawURL == "" {
		return nil, noopClose, nil
	}
	client, err := redisx.NewClient(ctx, rawURL)
	if err != nil {
		return nil, nil, err
	}
	return client, func() { _ = client.Close() }, nil
}

func closeAll(closers []func()) {
	for i := len(closers) - 1; i >= 0; i-- {
		if closers[i] != nil {
			closers[i]()
		}
	}
}
