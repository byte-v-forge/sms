package app

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/byte-v-forge/common-lib/eventoutbox"
	"github.com/byte-v-forge/sms/internal/core"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const smsPlatformEventOutboxTable = "sms_platform_event_outbox"

type PostgresOrderStore struct{ pool *pgxpool.Pool }

func NewPostgresOrderStore(ctx context.Context, dsn string) (*PostgresOrderStore, error) {
	pool, err := newRequiredPostgresPool(ctx, dsn)
	if err != nil {
		return nil, err
	}
	store := &PostgresOrderStore{pool: pool}
	if err := validatePostgresTables(ctx, pool, "sms_orders", smsPlatformEventOutboxTable); err != nil {
		pool.Close()
		return nil, err
	}
	return store, nil
}

func (s *PostgresOrderStore) Close() {
	if s != nil && s.pool != nil {
		s.pool.Close()
	}
}
func (s *PostgresOrderStore) Save(ctx context.Context, order core.Order, events ...eventoutbox.Record) error {
	return s.upsert(ctx, order, events...)
}
func (s *PostgresOrderStore) Update(ctx context.Context, order core.Order, events ...eventoutbox.Record) error {
	return s.upsert(ctx, order, events...)
}

func (s *PostgresOrderStore) upsert(ctx context.Context, order core.Order, events ...eventoutbox.Record) error {
	return s.withTx(ctx, func(tx pgx.Tx) error {
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
		if err != nil {
			return err
		}
		return insertOutboxRecords(ctx, tx, events)
	})
}

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

func (s *PostgresOrderStore) withTx(ctx context.Context, fn func(pgx.Tx) error) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func insertOutboxRecords(ctx context.Context, tx pgx.Tx, records []eventoutbox.Record) error {
	now := time.Now().Unix()
	for _, record := range records {
		if record.EventID == "" {
			continue
		}
		if err := eventoutbox.InsertRecordPgx(ctx, tx, smsPlatformEventOutboxTable, record, now); err != nil {
			return err
		}
	}
	return nil
}

func orderValues(order core.Order) []any {
	errorCode := errorCode(order.LastError)
	return []any{
		order.ID, order.RequestID, order.ProviderKey, order.UpstreamOrderID,
		order.Target.ApplicationKey, order.Target.CountryISO2, order.Target.CountryCallingCode,
		order.PhoneNumber.E164, string(order.Status), order.Price.CurrencyCode, order.Price.AmountDecimal,
		timeOrNil(order.AcquiredAt), timeOrNil(order.ExpiresAt), timeOrNil(order.UpdatedAt),
		timeOrNil(order.CancelAllowedAt), errorCode,
	}
}

func orderColumns() string {
	return `order_id, request_id, provider_key, upstream_order_id, target_application_key, target_country_iso2, target_country_calling_code, phone_e164, status, price_currency, price_amount, acquired_at, expires_at, updated_at, cancel_allowed_at, last_error_code`
}

func scanOrder(row pgx.Row) (core.Order, error) {
	var r orderRecord
	if err := row.Scan(&r.id, &r.requestID, &r.providerKey, &r.upstreamOrderID, &r.targetApplicationKey, &r.targetCountryISO2, &r.targetCallingCode, &r.phoneE164, &r.status, &r.priceCurrency, &r.priceAmount, &r.acquiredAt, &r.expiresAt, &r.updatedAt, &r.cancelAllowedAt, &r.lastErrorCode); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core.Order{}, core.NewError(core.CodeOrderNotFound, "order not found", false)
		}
		return core.Order{}, err
	}
	return r.toCore(), nil
}

type orderRecord struct {
	id                   string
	requestID            string
	providerKey          string
	upstreamOrderID        string
	targetApplicationKey string
	targetCountryISO2    string
	targetCallingCode    string
	phoneE164            string
	status               string
	priceCurrency        string
	priceAmount          string
	acquiredAt           sql.NullTime
	expiresAt            sql.NullTime
	updatedAt            time.Time
	cancelAllowedAt      sql.NullTime
	lastErrorCode        string
}

func (r orderRecord) toCore() core.Order {
	return core.Order{ID: r.id, RequestID: r.requestID, ProviderKey: r.providerKey, UpstreamOrderID: r.upstreamOrderID, Target: core.Target{ApplicationKey: r.targetApplicationKey, CountryISO2: r.targetCountryISO2, CountryCallingCode: r.targetCallingCode}, PhoneNumber: core.PhoneNumber{E164: r.phoneE164, CountryISO2: r.targetCountryISO2, CountryCallingCode: r.targetCallingCode}, Status: core.OrderStatus(r.status), Price: core.Money{CurrencyCode: r.priceCurrency, AmountDecimal: r.priceAmount}, AcquiredAt: nullableTime(r.acquiredAt), ExpiresAt: nullableTime(r.expiresAt), UpdatedAt: r.updatedAt, CancelAllowedAt: nullableTime(r.cancelAllowedAt), LastError: errorFromCode(r.lastErrorCode)}
}

func timeOrNil(value time.Time) any {
	if value.IsZero() {
		return nil
	}
	return value
}
func nullableTime(value sql.NullTime) time.Time {
	if !value.Valid {
		return time.Time{}
	}
	return value.Time
}
func errorCode(err *core.Error) string {
	if err == nil {
		return ""
	}
	return string(err.Code)
}
func errorFromCode(code string) *core.Error {
	if code == "" {
		return nil
	}
	return &core.Error{Code: core.ErrorCode(code)}
}
