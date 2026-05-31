package app

import (
	"context"
	"time"

	"github.com/byte-v-forge/common-lib/eventoutbox"
	"github.com/byte-v-forge/sms/internal/core"
	"github.com/jackc/pgx/v5"
)

func (s *PostgresOrderStore) RecordCode(ctx context.Context, order core.Order, code core.SMSCode, events ...eventoutbox.Record) error {
	return s.withTx(ctx, func(tx pgx.Tx) error {
		if err := upsertOrder(ctx, tx, order); err != nil {
			return err
		}
		receivedAt := code.ReceivedAt
		if receivedAt.IsZero() {
			receivedAt = time.Now().UTC()
		}
		if _, err := tx.Exec(ctx, `
INSERT INTO sms_order_codes (
  order_id, code_value, message_text, received_at
) VALUES ($1,$2,$3,$4)
ON CONFLICT (order_id, code_value, received_at) DO UPDATE SET
  message_text = EXCLUDED.message_text
`, order.ID, code.Value, code.MessageText, receivedAt); err != nil {
			return err
		}
		return insertOutboxRecords(ctx, tx, events)
	})
}

func (s *PostgresOrderStore) ListCodes(ctx context.Context, orderIDs []string, limitPerOrder int) ([]core.OrderCode, error) {
	if len(orderIDs) == 0 {
		return []core.OrderCode{}, nil
	}
	if limitPerOrder <= 0 {
		limitPerOrder = 10
	}
	if limitPerOrder > 50 {
		limitPerOrder = 50
	}
	rows, err := s.pool.Query(ctx, `
SELECT order_id, code_value, message_text, received_at
FROM (
  SELECT
    order_id,
    code_value,
    message_text,
    received_at,
    row_number() OVER (PARTITION BY order_id ORDER BY received_at DESC, code_value ASC) AS rn
  FROM sms_order_codes
  WHERE order_id = ANY($1::text[])
) codes
WHERE rn <= $2
ORDER BY order_id ASC, received_at DESC, code_value ASC
`, orderIDs, limitPerOrder)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []core.OrderCode{}
	for rows.Next() {
		var item core.OrderCode
		if err := rows.Scan(&item.OrderID, &item.Code.Value, &item.Code.MessageText, &item.Code.ReceivedAt); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}
