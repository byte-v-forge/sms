package app

import (
	"database/sql"
	"errors"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/jackc/pgx/v5"
)

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
	upstreamOrderID      string
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
