package fivesim

import (
	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/jsonx"
	"github.com/byte-v-forge/sms/internal/providers/phone"
)

func providerOrderFromPayload(payload order, request core.ProviderAcquireRequest, currencyCode string) (core.ProviderOrder, error) {
	orderID := jsonx.Scalar(payload.ID)
	if orderID == "" {
		return core.ProviderOrder{}, core.NewError(core.CodeUpstreamRejected, "missing 5sim order id", false)
	}
	e164, national := phone.Normalize(payload.Phone, request.Target.CountryISO2, request.Target.CountryCallingCode)
	return core.ProviderOrder{
		UpstreamOrderID: orderID,
		PhoneNumber: core.PhoneNumber{
			E164:               e164,
			NationalNumber:     national,
			CountryISO2:        request.Target.CountryISO2,
			CountryCallingCode: request.Target.CountryCallingCode,
		},
		Price:      core.Money{CurrencyCode: currencyCode, AmountDecimal: jsonx.Scalar(payload.Price)},
		AcquiredAt: parseTime(payload.CreatedAt),
		ExpiresAt:  parseTime(payload.Expires),
	}, nil
}
