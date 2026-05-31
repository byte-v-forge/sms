package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *PostgresOrderStore) Get(ctx context.Context, orderID string) (core.Order, error) {
	row := s.pool.QueryRow(ctx, `SELECT `+orderColumns()+` FROM sms_orders WHERE order_id = $1`, orderID)
	return scanOrder(row)
}

func (s *PostgresOrderStore) List(ctx context.Context, includeFinal bool, limit int) ([]core.Order, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 500 {
		limit = 500
	}
	rows, err := s.pool.Query(ctx, `
SELECT `+orderColumns()+`
FROM sms_orders
WHERE $1 OR status NOT IN ('completed', 'canceled', 'expired', 'failed')
ORDER BY updated_at DESC, order_id ASC
LIMIT $2
`, includeFinal, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []core.Order{}
	for rows.Next() {
		order, err := scanOrder(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, order)
	}
	return out, rows.Err()
}
