package smsbower

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/providers/handlerapi"
	"github.com/byte-v-forge/sms/internal/providers/phone"
)

func (c *Client) AcquireNumber(ctx context.Context, request core.ProviderAcquireRequest) (core.ProviderOrder, error) {
	if strings.TrimSpace(request.Route.UpstreamServiceKey) == "" {
		return core.ProviderOrder{}, core.NewError(core.CodeValidationFailed, "smsbower service is required", false)
	}
	if strings.TrimSpace(request.Route.ProviderCountryID) == "" {
		return core.ProviderOrder{}, core.NewError(core.CodeValidationFailed, "smsbower country is required", false)
	}
	if strings.TrimSpace(request.Route.UpstreamProviderID) == "" {
		return core.ProviderOrder{}, core.NewError(core.CodeValidationFailed, "smsbower upstream provider id is required", false)
	}

	params := c.acquireParams(request)
	result, err := c.api.Do(ctx, "getNumberV2", params)
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
	return c.api.SetStatus(ctx, upstreamOrderID, action, statusForAction)
}

func (c *Client) GetBalance(ctx context.Context) (core.Money, error) {
	return c.api.GetBalance(ctx)
}
func (c *Client) acquireParams(request core.ProviderAcquireRequest) url.Values {
	params := url.Values{}
	params.Set("service", request.Route.UpstreamServiceKey)
	params.Set("country", request.Route.ProviderCountryID)
	params.Set("providerIds", strings.TrimSpace(request.Route.UpstreamProviderID))
	return params
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
	orderID := rawJSONScalar(payload.OrderID)
	rawPhone := rawJSONScalar(payload.PhoneNumber)
	if orderID == "" || rawPhone == "" {
		return core.ProviderOrder{}, core.NewError(core.CodeUpstreamRejected, "missing activationId or phoneNumber in getNumberV2 response", false)
	}
	e164, national := phone.Normalize(rawPhone, request.Target.CountryISO2, request.Target.CountryCallingCode)
	return core.ProviderOrder{
		UpstreamOrderID:          orderID,
		PhoneNumber:              core.PhoneNumber{E164: e164, NationalNumber: national, CountryISO2: request.Target.CountryISO2, CountryCallingCode: request.Target.CountryCallingCode},
		Price:                    core.Money{AmountDecimal: rawJSONScalar(payload.OrderCost)},
		AcquiredAt:               parseOrderTimeText(rawJSONScalar(payload.OrderTime)),
		CanRequestAdditionalCode: providerTruthy(payload.CanGetAnotherSMS),
	}, nil
}

func providerTruthy(raw json.RawMessage) bool {
	switch strings.ToLower(rawJSONScalar(raw)) {
	case "1", "true", "yes":
		return true
	default:
		return false
	}
}
