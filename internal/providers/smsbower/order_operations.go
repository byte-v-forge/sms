package smsbower

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/jsonx"
	"github.com/byte-v-forge/sms/internal/providers/handlerapi"
	"github.com/byte-v-forge/sms/internal/providers/phone"
)

func (c *Client) AcquireNumber(ctx context.Context, request core.ProviderAcquireRequest) (core.ProviderOrder, error) {
	result, err := c.api.GetNumberV2(ctx, request, handlerapi.GetNumberV2Config{
		ProviderName:       "smsbower",
		ProviderIDParam:    "providerIds",
		ProviderIDRequired: true,
	})
	if err != nil {
		return core.ProviderOrder{}, err
	}
	order, err := c.parseGetNumberV2(result, request)
	if err == nil {
		return order, nil
	}
	if isProviderTextError(result) {
		return core.ProviderOrder{}, handlerapi.MapTextError(result)
	}
	return core.ProviderOrder{}, err
}

func (c *Client) GetStatus(ctx context.Context, upstreamOrderID string) (core.ProviderCodeResult, error) {
	return c.api.GetStatus(ctx, upstreamOrderID, parseStatus)
}

func (c *Client) SetStatus(ctx context.Context, upstreamOrderID string, action core.ProviderAction) error {
	return c.api.SetActivationStatus(ctx, upstreamOrderID, action, "smsbower")
}

func (c *Client) GetBalance(ctx context.Context) (core.Money, error) {
	return c.api.GetBalance(ctx)
}

func (c *Client) parseGetNumberV2(result string, request core.ProviderAcquireRequest) (core.ProviderOrder, error) {
	var payload struct {
		OrderID          json.RawMessage `json:"activationId"`
		PhoneNumber      json.RawMessage `json:"phoneNumber"`
		OrderCost        json.RawMessage `json:"activationCost"`
		CanGetAnotherSMS json.RawMessage `json:"canGetAnotherSms"`
		OrderTime        json.RawMessage `json:"activationTime"`
	}
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

func providerTruthy(raw json.RawMessage) bool {
	switch strings.ToLower(jsonx.Scalar(raw)) {
	case "1", "true", "yes":
		return true
	default:
		return false
	}
}
