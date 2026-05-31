package app

import (
	"github.com/byte-v-forge/common-lib/eventbus"
	"github.com/byte-v-forge/common-lib/eventoutbox"
)

func NewOrderOutboxProcessor(store *PostgresOrderStore, publisher eventbus.Publisher) eventoutbox.PendingProcessor {
	if store == nil {
		return nil
	}
	return eventoutbox.NewPgxProcessor(store.pool, smsPlatformEventOutboxTable, publisher)
}
