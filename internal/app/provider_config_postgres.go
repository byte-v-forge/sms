package app

import (
	"context"

	providerspi "github.com/byte-v-forge/sms/internal/providers/spi"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresProviderConfigStore struct {
	pool      *pgxpool.Pool
	providers *providerspi.Registry
}

func NewPostgresProviderConfigStore(ctx context.Context, dsn string, providers *providerspi.Registry) (*PostgresProviderConfigStore, error) {
	pool, err := newPostgresPool(ctx, dsn)
	if err != nil {
		return nil, err
	}
	store := &PostgresProviderConfigStore{pool: pool, providers: providers}
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
