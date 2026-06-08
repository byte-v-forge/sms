package eventoutbox

import (
	"context"
	"fmt"
	"time"
)

func InsertRecordPgx(ctx context.Context, tx PgxTx, table string, record Record, now int64) error {
	if tx == nil {
		return ErrNilDB
	}
	tableName, err := postgresIdentifier(table)
	if err != nil {
		return err
	}
	if now <= 0 {
		now = time.Now().Unix()
	}
	_, err = tx.Exec(ctx, fmt.Sprintf(`
		INSERT INTO %s (
			event_id, subject, event_name, idempotency_key, envelope, status,
			attempt_count, next_attempt_at, last_error, published_at, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, 0, 0, '', 0, $7, $7)
		ON CONFLICT (event_id) DO NOTHING
	`, tableName), record.EventID, record.Subject, record.EventName, record.IdempotencyKey, record.Envelope, StatusPending, now)
	return err
}

func ClaimPendingPgx(ctx context.Context, tx PgxTx, table string, batch int, now int64) ([]Row, error) {
	if tx == nil {
		return nil, ErrNilDB
	}
	tableName, err := postgresIdentifier(table)
	if err != nil {
		return nil, err
	}
	if batch <= 0 {
		batch = DefaultBatch
	}
	if now <= 0 {
		now = time.Now().Unix()
	}
	rows, err := tx.Query(ctx, fmt.Sprintf(`
		SELECT event_id, envelope, attempt_count
		FROM %s
		WHERE status = $1 AND next_attempt_at <= $2
		ORDER BY created_at ASC
		LIMIT $3
		FOR UPDATE SKIP LOCKED
	`, tableName), StatusPending, now, batch)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []Row{}
	for rows.Next() {
		var row Row
		if err := rows.Scan(&row.EventID, &row.Envelope, &row.AttemptCount); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, rows.Err()
}
