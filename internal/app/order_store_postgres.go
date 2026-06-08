package app

import (
	"context"
	"errors"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/eventoutbox"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const smsPlatformEventOutboxTable = "sms_platform_event_outbox"

type PostgresOrderStore struct{ pool *pgxpool.Pool }

func NewPostgresOrderStore(ctx context.Context, pool *pgxpool.Pool) (*PostgresOrderStore, error) {
	if pool == nil {
		return nil, errors.New("postgres order store pool is required")
	}
	store := &PostgresOrderStore{pool: pool}
	if err := validatePostgresTables(ctx, pool, "sms_orders", "sms_order_codes", smsPlatformEventOutboxTable); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *PostgresOrderStore) Save(ctx context.Context, order core.Order, events ...eventoutbox.Record) error {
	return s.upsert(ctx, order, events...)
}

func (s *PostgresOrderStore) Update(ctx context.Context, order core.Order, events ...eventoutbox.Record) error {
	return s.upsert(ctx, order, events...)
}

func (s *PostgresOrderStore) upsert(ctx context.Context, order core.Order, events ...eventoutbox.Record) error {
	return s.withTx(ctx, func(tx pgx.Tx) error {
		if err := upsertOrder(ctx, tx, order); err != nil {
			return err
		}
		return insertOutboxRecords(ctx, tx, events)
	})
}
