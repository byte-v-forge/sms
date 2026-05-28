package app

import (
	"context"

	"github.com/byte-v-forge/common-lib/eventbus"
	"github.com/byte-v-forge/common-lib/eventoutbox"
)

type OrderOutboxProcessor struct {
	store     *PostgresOrderStore
	publisher eventbus.Publisher
}

func NewOrderOutboxProcessor(store *PostgresOrderStore, publisher eventbus.Publisher) *OrderOutboxProcessor {
	return &OrderOutboxProcessor{store: store, publisher: publisher}
}

func (p *OrderOutboxProcessor) PublishPending(ctx context.Context, batch int) (int, error) {
	if p == nil || p.store == nil || p.publisher == nil {
		return 0, nil
	}
	tx, err := p.store.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	rows, err := eventoutbox.ClaimPendingPgx(ctx, tx, smsPlatformEventOutboxTable, batch, 0)
	if err != nil {
		return 0, err
	}
	updates, err := eventoutbox.NewPgxUpdates(tx, smsPlatformEventOutboxTable)
	if err != nil {
		return 0, err
	}
	published, err := eventoutbox.PublishRows(ctx, p.publisher, rows, updates, eventoutbox.PublishOptions{})
	if err != nil {
		return published, err
	}
	return published, tx.Commit(ctx)
}
