package app

import (
	"context"
	"errors"

	providerspi "github.com/byte-v-forge/sms/internal/providers/spi"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresProviderConfigStore struct {
	pool      *pgxpool.Pool
	providers *providerspi.Registry
}

func NewPostgresProviderConfigStore(ctx context.Context, pool *pgxpool.Pool, providers *providerspi.Registry) (*PostgresProviderConfigStore, error) {
	if pool == nil {
		return nil, errors.New("postgres provider config pool is required")
	}
	store := &PostgresProviderConfigStore{pool: pool, providers: providers}
	if err := validatePostgresTables(ctx, pool, "sms_provider_configs"); err != nil {
		return nil, err
	}
	return store, nil
}
