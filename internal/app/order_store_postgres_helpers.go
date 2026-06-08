package app

import (
	"github.com/byte-v-forge/sms/internal/core"
)

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
