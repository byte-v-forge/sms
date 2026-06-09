package main

import (
	"context"
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/app"
	"github.com/byte-v-forge/sms/internal/platform/redisx"
	"github.com/redis/go-redis/v9"
)

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

func configureRuntimeCacheStores(store *runtimeStores, redisClient *redis.Client, history orderHistoryStore, clock app.SystemClock) {
	if redisClient == nil {
		store.codeSecrets = app.NewSMSCodeSecretStore(app.NewMemoryTTLStringStore("sms:code-secret", 30*time.Minute, clock), clock)
		store.routeHealth = app.NewMemoryRouteHealthStore(clock)
		return
	}
	activeStore := app.NewRedisOrderStore(redisx.NewStringStore(redisClient, "sms:order", 30*time.Minute), clock)
	store.orders = app.NewCompositeOrderStore(activeStore, history)
	store.codeSecrets = app.NewSMSCodeSecretStore(redisx.NewStringStore(redisClient, "sms:code-secret", 30*time.Minute), clock)
	store.routeHealth = app.NewRedisRouteHealthStore(redisClient)
}
