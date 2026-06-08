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
