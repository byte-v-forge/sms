package app

import (
	"context"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/eventoutbox"
	"github.com/byte-v-forge/sms/internal/platform/secretref"
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
		if !secretref.Configured(code.SecretRef) {
			return core.NewError(core.CodeInternal, "sms code secret ref is required", false)
		}
		expiresAt := receivedAt.Add(30 * time.Minute)
		if ts := code.SecretRef.GetExpiresAt(); ts != nil {
			expiresAt = ts.AsTime()
		}
		if _, err := tx.Exec(ctx, `
INSERT INTO sms_order_codes (
  order_id, code_secret_id, message_text, received_at, expires_at
) VALUES ($1,$2,$3,$4,$5)
ON CONFLICT (order_id, code_secret_id, received_at) DO UPDATE SET
  message_text = EXCLUDED.message_text,
  expires_at = EXCLUDED.expires_at
`, order.ID, code.SecretRef.GetSecretId(), code.MessageText, receivedAt, expiresAt); err != nil {
			return err
		}
		return insertOutboxRecords(ctx, tx, events)
	})
}

func (s *PostgresOrderStore) CodeSecretExists(ctx context.Context, orderID string, secretID string) (bool, error) {
	var exists bool
	err := s.pool.QueryRow(ctx, `
SELECT EXISTS (
  SELECT 1
  FROM sms_order_codes
  WHERE order_id = $1 AND code_secret_id = $2
)
`, orderID, secretID).Scan(&exists)
	return exists, err
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
	out := []core.OrderCode{}
	for rows.Next() {
		var item core.OrderCode
		var secretID string
		var expiresAt time.Time
		if err := rows.Scan(&item.OrderID, &secretID, &item.Code.MessageText, &item.Code.ReceivedAt, &expiresAt); err != nil {
			return nil, err
		}
		item.Code.SecretRef = secretref.New("sms", "sms_otp", secretID, expiresAt)
		out = append(out, item)
	}
	return out, rows.Err()
}
