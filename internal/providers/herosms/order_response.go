package herosms

import (
	"encoding/json"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/jsonx"
	"github.com/byte-v-forge/sms/internal/providers/phone"
)

type heroSMSGetNumberV2Response struct {
	ActivationID         json.RawMessage `json:"activationId"`
	PhoneNumber          string          `json:"phoneNumber"`
	ActivationCost       json.RawMessage `json:"activationCost"`
	Currency             json.RawMessage `json:"currency"`
	CanGetAnotherSMS     json.RawMessage `json:"canGetAnotherSms"`
	ActivationTime       string          `json:"activationTime"`
	ActivationEndTime    string          `json:"activationEndTime"`
	ActivationOperator   string          `json:"activationOperator"`
	ServiceCode          string          `json:"serviceCode"`
	CountryPhoneCode     json.RawMessage `json:"countryPhoneCode"`
	ActivationStatusCode json.RawMessage `json:"status"`
}

func heroSMSProviderOrder(orderID string, rawPhone string, payload heroSMSGetNumberV2Response, request core.ProviderAcquireRequest) core.ProviderOrder {
	e164, national := phone.Normalize(rawPhone, request.Target.CountryISO2, request.Target.CountryCallingCode)
	return core.ProviderOrder{
		UpstreamOrderID: orderID,
		PhoneNumber: core.PhoneNumber{
			E164:               e164,
			NationalNumber:     national,
			CountryISO2:        request.Target.CountryISO2,
			CountryCallingCode: request.Target.CountryCallingCode,
		},
		Price:                    core.Money{CurrencyCode: heroSMSCurrencyCode(payload.Currency), AmountDecimal: jsonx.FirstScalar(payload.ActivationCost)},
		AcquiredAt:               parseHeroSMSTime(payload.ActivationTime),
		ExpiresAt:                parseHeroSMSTime(payload.ActivationEndTime),
		CanRequestAdditionalCode: heroSMSBool(payload.CanGetAnotherSMS),
	}
}
