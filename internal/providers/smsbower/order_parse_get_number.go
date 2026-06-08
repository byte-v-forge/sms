package smsbower

import (
	"encoding/json"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/jsonx"
	"github.com/byte-v-forge/sms/internal/providers/phone"
)

func (c *Client) parseGetNumberV2(result string, request core.ProviderAcquireRequest) (core.ProviderOrder, error) {
	var payload getNumberV2Payload
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		return core.ProviderOrder{}, core.NewError(core.CodeUpstreamRejected, "bad getNumberV2 json response", false)
	}
	orderID := jsonx.Scalar(payload.OrderID)
	rawPhone := jsonx.Scalar(payload.PhoneNumber)
	if orderID == "" || rawPhone == "" {
		return core.ProviderOrder{}, core.NewError(core.CodeUpstreamRejected, "missing activationId or phoneNumber in getNumberV2 response", false)
	}
	e164, national := phone.Normalize(rawPhone, request.Target.CountryISO2, request.Target.CountryCallingCode)
	return core.ProviderOrder{
		UpstreamOrderID:          orderID,
		PhoneNumber:              core.PhoneNumber{E164: e164, NationalNumber: national, CountryISO2: request.Target.CountryISO2, CountryCallingCode: request.Target.CountryCallingCode},
		Price:                    core.Money{AmountDecimal: jsonx.Scalar(payload.OrderCost)},
		AcquiredAt:               parseOrderTimeText(jsonx.Scalar(payload.OrderTime)),
		CanRequestAdditionalCode: providerTruthy(payload.CanGetAnotherSMS),
	}, nil
}
