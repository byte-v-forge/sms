package eventoutbox

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func ClaimPendingPgx(ctx context.Context, tx PgxTx, table string, batch int, now int64) ([]Row, error) {
	if tx == nil {
		return nil, ErrNilDB
	}
	tableName, err := postgresIdentifier(table)
	if err != nil {
		return nil, err
	}
	rows, err := tx.Query(ctx, fmt.Sprintf(`
		SELECT event_id, envelope, attempt_count
		FROM %s
		WHERE status = $1 AND next_attempt_at <= $2
		ORDER BY created_at ASC
		LIMIT $3
		FOR UPDATE SKIP LOCKED
	`, tableName), StatusPending, resolveUnixTime(now), resolveBatchSize(batch))
	if err != nil {
		return nil, err
	}
	return scanPendingRows(rows)
}

func scanPendingRows(rows pgx.Rows) ([]Row, error) {
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
