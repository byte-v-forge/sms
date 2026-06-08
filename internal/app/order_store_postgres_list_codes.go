package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *PostgresOrderStore) ListCodes(ctx context.Context, orderIDs []string, limitPerOrder int) ([]core.OrderCode, error) {
	if len(orderIDs) == 0 {
		return []core.OrderCode{}, nil
	}
	limitPerOrder = normalizePostgresOrderCodeLimit(limitPerOrder)
	rows, err := s.pool.Query(ctx, `
SELECT order_id, code_secret_id, message_text, received_at, expires_at
FROM (
  SELECT
    order_id,
    code_secret_id,
    message_text,
    received_at,
    expires_at,
    row_number() OVER (PARTITION BY order_id ORDER BY received_at DESC, code_secret_id ASC) AS rn
  FROM sms_order_codes
  WHERE order_id = ANY($1::text[])
) codes
WHERE rn <= $2
ORDER BY order_id ASC, received_at DESC, code_secret_id ASC
`, orderIDs, limitPerOrder)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanPostgresOrderCodes(rows)
}

func normalizePostgresOrderCodeLimit(limitPerOrder int) int {
	if limitPerOrder <= 0 {
		return 10
	}
	if limitPerOrder > 50 {
		return 50
	}
	return limitPerOrder
}
