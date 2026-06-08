package fivesim

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/providers/phone"
)

func (c *Client) AcquireNumber(ctx context.Context, request core.ProviderAcquireRequest) (core.ProviderOrder, error) {
	country := strings.TrimSpace(request.Route.ProviderCountryID)
	if country == "" {
		return core.ProviderOrder{}, core.NewError(core.CodeValidationFailed, "5sim country is required", false)
	}
	product := strings.TrimSpace(request.Route.UpstreamServiceKey)
	if product == "" {
		return core.ProviderOrder{}, core.NewError(core.CodeValidationFailed, "5sim product is required", false)
	}
	operator, err := operatorForRoute(request.Route)
	if err != nil {
		return core.ProviderOrder{}, err
	}
	params := c.buyParams()
	path := fmt.Sprintf("/v1/user/buy/activation/%s/%s/%s", url.PathEscape(country), url.PathEscape(operator), url.PathEscape(product))

	var payload order
	if err := c.getJSON(ctx, path, params, true, &payload); err != nil {
		return core.ProviderOrder{}, err
	}
	orderID := rawJSONScalar(payload.ID)
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
		Price:      core.Money{CurrencyCode: c.currencyCode, AmountDecimal: rawJSONScalar(payload.Price)},
		AcquiredAt: parseTime(payload.CreatedAt),
		ExpiresAt:  parseTime(payload.Expires),
	}, nil
}

func (c *Client) GetStatus(ctx context.Context, upstreamOrderID string) (core.ProviderCodeResult, error) {
	var payload order
	if err := c.getJSON(ctx, "/v1/user/check/"+url.PathEscape(upstreamOrderID), nil, true, &payload); err != nil {
		return core.ProviderCodeResult{}, err
	}
	return orderToCodeResult(payload), nil
}

func (c *Client) SetStatus(ctx context.Context, upstreamOrderID string, action core.ProviderAction) error {
	switch action {
	case core.ActionMarkMessageSent:
		return core.NewError(core.CodeUnsupportedOperation, "5sim does not expose mark-message-sent status", false)
	case core.ActionRequestAdditional:
		return nil
	case core.ActionCompleteOrder:
		var payload order
		return c.getJSON(ctx, "/v1/user/finish/"+url.PathEscape(upstreamOrderID), nil, true, &payload)
	case core.ActionCancelOrder:
		var payload order
		return c.getJSON(ctx, "/v1/user/cancel/"+url.PathEscape(upstreamOrderID), nil, true, &payload)
	default:
		return core.NewError(core.CodeUnsupportedOperation, "unsupported 5sim status action", false)
	}
}

func (c *Client) GetBalance(ctx context.Context) (core.Money, error) {
	var payload struct {
		Balance json.RawMessage `json:"balance"`
	}
	if err := c.getJSON(ctx, "/v1/user/profile", nil, true, &payload); err != nil {
		return core.Money{}, err
	}
	return core.Money{CurrencyCode: c.currencyCode, AmountDecimal: rawJSONScalar(payload.Balance)}, nil
}

func (c *Client) buyParams() url.Values {
	params := url.Values{}
	setBool(params, "reuse", c.reuse)
	setBool(params, "voice", c.voice)
	setBool(params, "forwarding", c.forwarding)
	if c.forwardingNumber != "" {
		params.Set("number", c.forwardingNumber)
	}
	if c.ref != "" {
		params.Set("ref", c.ref)
	}
	return params
}

func operatorForRoute(route core.Route) (string, error) {
	operator := strings.TrimSpace(route.UpstreamProviderID)
	if operator == "" {
		return "", core.NewError(core.CodeValidationFailed, "5sim operator is required", false)
	}
	return operator, nil
}
