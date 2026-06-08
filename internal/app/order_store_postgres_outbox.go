package app

import (
	"context"
	"time"

	"github.com/byte-v-forge/sms/internal/platform/eventoutbox"
	"github.com/jackc/pgx/v5"
)

func insertOutboxRecords(ctx context.Context, tx pgx.Tx, records []eventoutbox.Record) error {
	now := time.Now().Unix()
	for _, record := range records {
		if record.EventID == "" {
			continue
		}
		if err := eventoutbox.InsertRecordPgx(ctx, tx, smsPlatformEventOutboxTable, record, now); err != nil {
			return err
		}
	}
	return nil
}
