package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/eventoutbox"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const smsPlatformEventOutboxTable = "sms_platform_event_outbox"

type PostgresOrderStore struct{ pool *pgxpool.Pool }

func NewPostgresOrderStore(ctx context.Context, dsn string) (*PostgresOrderStore, error) {
	pool, err := newPostgresPool(ctx, dsn)
	if err != nil {
		return nil, err
	}
	store := &PostgresOrderStore{pool: pool}
	if err := validatePostgresTables(ctx, pool, "sms_orders", "sms_order_codes", smsPlatformEventOutboxTable); err != nil {
		pool.Close()
		return nil, err
	}
	return store, nil
}

func (s *PostgresOrderStore) Close() {
	if s != nil && s.pool != nil {
		s.pool.Close()
	}
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
