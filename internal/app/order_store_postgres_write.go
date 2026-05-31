package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/jackc/pgx/v5"
)

func upsertOrder(ctx context.Context, tx pgx.Tx, order core.Order) error {
	_, err := tx.Exec(ctx, `
INSERT INTO sms_orders (
  order_id, request_id, provider_key, upstream_order_id,
  target_application_key, target_country_iso2, target_country_calling_code,
  phone_e164, status, price_currency, price_amount, acquired_at, expires_at, updated_at,
  cancel_allowed_at, last_error_code
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
ON CONFLICT (order_id) DO UPDATE SET
  request_id = EXCLUDED.request_id,
  provider_key = EXCLUDED.provider_key,
  upstream_order_id = EXCLUDED.upstream_order_id,
  target_application_key = EXCLUDED.target_application_key,
  target_country_iso2 = EXCLUDED.target_country_iso2,
  target_country_calling_code = EXCLUDED.target_country_calling_code,
  phone_e164 = EXCLUDED.phone_e164,
  status = EXCLUDED.status,
  price_currency = EXCLUDED.price_currency,
  price_amount = EXCLUDED.price_amount,
  acquired_at = EXCLUDED.acquired_at,
  expires_at = EXCLUDED.expires_at,
  updated_at = EXCLUDED.updated_at,
  cancel_allowed_at = EXCLUDED.cancel_allowed_at,
  last_error_code = EXCLUDED.last_error_code
`, orderValues(order)...)
	return err
}
