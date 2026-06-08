package eventoutbox

import (
	"context"
	"fmt"
)

func NewPgxUpdates(tx PgxTx, table string) (Updates, error) {
	if tx == nil {
		return nil, ErrNilDB
	}
	tableName, err := postgresIdentifier(table)
	if err != nil {
		return nil, err
	}
	return pgxUpdates{tx: tx, table: tableName}, nil
}

type pgxUpdates struct {
	tx    PgxTx
	table string
}

func (u pgxUpdates) MarkPublished(ctx context.Context, eventID string, publishedAt int64) error {
	_, err := u.tx.Exec(ctx, fmt.Sprintf(`
		UPDATE %s
		SET status = $1, published_at = $2, updated_at = $2, last_error = ''
		WHERE event_id = $3
	`, u.table), StatusPublished, publishedAt, eventID)
	return err
}

func (u pgxUpdates) MarkRetry(ctx context.Context, eventID string, attemptCount int32, nextAttemptAt int64, lastError string, updatedAt int64) error {
	_, err := u.tx.Exec(ctx, fmt.Sprintf(`
		UPDATE %s
		SET attempt_count = $1, next_attempt_at = $2, last_error = $3, updated_at = $4
		WHERE event_id = $5
	`, u.table), attemptCount, nextAttemptAt, lastError, updatedAt, eventID)
	return err
}

func (u pgxUpdates) MarkDiscarded(ctx context.Context, eventID string, lastError string, updatedAt int64) error {
	_, err := u.tx.Exec(ctx, fmt.Sprintf(`
		UPDATE %s
		SET status = $1, last_error = $2, updated_at = $3
		WHERE event_id = $4
	`, u.table), StatusDiscarded, lastError, updatedAt, eventID)
	return err
}
