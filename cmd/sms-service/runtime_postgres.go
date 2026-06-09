package main

import (
	"context"
	"strings"

	"github.com/byte-v-forge/sms/internal/app"
	providerspi "github.com/byte-v-forge/sms/internal/providers/spi"
	"github.com/jackc/pgx/v5/pgxpool"
)

type orderHistoryStore interface {
	app.OrderStore
	app.OrderListStore
}

func newOptionalPostgresPool(ctx context.Context, cfg config) (*pgxpool.Pool, func(), error) {
	if strings.TrimSpace(cfg.PGDSN) == "" {
		return nil, noopClose, nil
	}
	pool, err := app.NewPostgresPool(ctx, cfg.PGDSN)
	if err != nil {
		return nil, nil, err
	}
	return pool, pool.Close, nil
}

func newProviderConfigStore(ctx context.Context, postgresPool *pgxpool.Pool, providers *providerspi.Registry, clock app.SystemClock) (app.ProviderConfigStore, error) {
	if postgresPool == nil {
		return app.NewMemoryProviderConfigStore(providers, clock), nil
	}
	store, err := app.NewPostgresProviderConfigStore(ctx, postgresPool, providers)
	if err != nil {
		return nil, err
	}
	return store, nil
}

func newOrderHistoryStore(ctx context.Context, postgresPool *pgxpool.Pool) (orderHistoryStore, *app.PostgresOrderStore, error) {
	if postgresPool == nil {
		store := app.NewMemoryOrderStore()
		return store, nil, nil
	}
	store, err := app.NewPostgresOrderStore(ctx, postgresPool)
	if err != nil {
		return nil, nil, err
	}
	return store, store, nil
}
