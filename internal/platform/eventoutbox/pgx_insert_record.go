package eventoutbox

import (
	"context"
	"fmt"
)

func InsertRecordPgx(ctx context.Context, tx PgxTx, table string, record Record, now int64) error {
	if tx == nil {
		return ErrNilDB
	}
	tableName, err := postgresIdentifier(table)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, fmt.Sprintf(`
		INSERT INTO %s (
			event_id, subject, event_name, idempotency_key, envelope, status,
			attempt_count, next_attempt_at, last_error, published_at, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, 0, 0, '', 0, $7, $7)
		ON CONFLICT (event_id) DO NOTHING
	`, tableName), record.EventID, record.Subject, record.EventName, record.IdempotencyKey, record.Envelope, StatusPending, resolveUnixTime(now))
	return err
}
