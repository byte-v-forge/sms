package app

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresProviderConfigStore struct{ pool *pgxpool.Pool }

func NewPostgresProviderConfigStore(ctx context.Context, dsn string) (*PostgresProviderConfigStore, error) {
	pool, err := newRequiredPostgresPool(ctx, dsn)
	if err != nil {
		return nil, err
	}
	store := &PostgresProviderConfigStore{pool: pool}
	if err := validatePostgresTables(ctx, pool, "sms_provider_configs"); err != nil {
		pool.Close()
		return nil, err
	}
	return store, nil
}

func (s *PostgresProviderConfigStore) Close() {
	if s != nil && s.pool != nil {
		s.pool.Close()
	}
}
