package smsbower

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
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
	params := url.Values{}
	params.Set("id", upstreamOrderID)
	result, err := c.api.Do(ctx, "getStatus", params)
	if err != nil {
		return core.ProviderCodeResult{}, err
	}
	return parseStatus(result)
}

func (c *Client) SetStatus(ctx context.Context, upstreamOrderID string, action core.ProviderAction) error {
	status, expected, err := statusForAction(action)
	if err != nil {
		return err
	}
	params := url.Values{}
	params.Set("id", upstreamOrderID)
	params.Set("status", status)
	result, err := c.api.Do(ctx, "setStatus", params)
	if err != nil {
		return err
	}
	if result != expected {
		return handlerapi.MapTextError(result)
	}
	return nil
}

func (c *Client) GetBalance(ctx context.Context) (core.Money, error) {
	result, err := c.api.Do(ctx, "getBalance", nil)
	if err != nil {
		return core.Money{}, err
	}
	const prefix = "ACCESS_BALANCE:"
	if !strings.HasPrefix(result, prefix) {
		return core.Money{}, handlerapi.MapTextError(result)
	}
	return core.Money{AmountDecimal: strings.TrimPrefix(result, prefix)}, nil
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
		OrderID          int64  `json:"activationId"`
		PhoneNumber      string `json:"phoneNumber"`
		OrderCost        string `json:"activationCost"`
		CountryCode      string `json:"countryCode"`
		CanGetAnotherSMS string `json:"canGetAnotherSms"`
		OrderTime        string `json:"activationTime"`
	}
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		return core.ProviderOrder{}, core.NewError(core.CodeUpstreamRejected, "bad getNumberV2 json response", false)
	}
	if payload.OrderID <= 0 {
		return core.ProviderOrder{}, core.NewError(core.CodeUpstreamRejected, "missing activationId in getNumberV2 response", false)
	}
	orderID := strconv.FormatInt(payload.OrderID, 10)
	e164, national := phone.Normalize(payload.PhoneNumber, request.Target.CountryISO2, request.Target.CountryCallingCode)
	return core.ProviderOrder{
		UpstreamOrderID:          orderID,
		PhoneNumber:              core.PhoneNumber{E164: e164, NationalNumber: national, CountryISO2: request.Target.CountryISO2, CountryCallingCode: request.Target.CountryCallingCode},
		Price:                    core.Money{AmountDecimal: strings.TrimSpace(payload.OrderCost)},
		AcquiredAt:               parseOrderTimeText(payload.OrderTime),
		CanRequestAdditionalCode: payload.CanGetAnotherSMS == "1",
	}, nil
}
