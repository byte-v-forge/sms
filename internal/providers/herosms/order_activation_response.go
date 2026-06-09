package herosms

import (
	"encoding/json"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/jsonx"
)

type activationPurchaseEnvelope struct {
	Data json.RawMessage `json:"data"`
}

func decodeActivationPurchase(raw []byte) (heroSMSGetNumberV2Response, error) {
	var envelope activationPurchaseEnvelope
	if err := json.Unmarshal(raw, &envelope); err == nil && len(envelope.Data) > 0 {
		raw = envelope.Data
	}
	var payload heroSMSGetNumberV2Response
	if err := json.Unmarshal(raw, &payload); err != nil {
		return heroSMSGetNumberV2Response{}, core.NewError(core.CodeUpstreamRejected, "bad hero sms activation response", false)
	}
	return payload, nil
}

func heroSMSProviderOrderFromPayload(payload heroSMSGetNumberV2Response, request core.ProviderAcquireRequest) (core.ProviderOrder, error) {
	orderID := jsonx.FirstScalar(payload.ActivationID)
	if orderID == "" || payload.PhoneNumber == "" {
		return core.ProviderOrder{}, core.NewError(core.CodeUpstreamRejected, "bad hero sms activation response", false)
	}
	return heroSMSProviderOrder(orderID, payload.PhoneNumber, payload, request), nil
}
