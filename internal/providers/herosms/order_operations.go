package herosms

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/providers/handlerapi"
	"github.com/byte-v-forge/sms/internal/providers/phone"
)

func (c *Client) AcquireNumber(ctx context.Context, request core.ProviderAcquireRequest) (core.ProviderOrder, error) {
	service := strings.TrimSpace(request.Route.UpstreamServiceKey)
	if service == "" {
		return core.ProviderOrder{}, core.NewError(core.CodeValidationFailed, "hero sms service is required", false)
	}
	country := strings.TrimSpace(request.Route.ProviderCountryID)
	if country == "" {
		return core.ProviderOrder{}, core.NewError(core.CodeValidationFailed, "hero sms provider country id is required", false)
	}
	params := url.Values{}
	params.Set("service", service)
	params.Set("country", country)
	if operator := strings.TrimSpace(request.Route.UpstreamProviderID); operator != "" {
		params.Set("operator", operator)
	}
	if maxPrice := strings.TrimSpace(request.Route.MaxPrice.AmountDecimal); maxPrice != "" {
		params.Set("maxPrice", maxPrice)
	}

	result, err := c.api.Do(ctx, "getNumberV2", params)
	if err != nil {
		return core.ProviderOrder{}, err
	}
	var payload heroSMSGetNumberV2Response
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		if orderID, rawPhone, ok := parseAccessNumber(result); ok {
			return heroSMSProviderOrder(orderID, rawPhone, heroSMSGetNumberV2Response{}, request), nil
		}
		return core.ProviderOrder{}, handlerapi.MapTextError(result)
	}
	orderID := firstHeroSMSScalar(payload.ActivationID)
	if orderID == "" || strings.TrimSpace(payload.PhoneNumber) == "" {
		return core.ProviderOrder{}, core.NewError(core.CodeUpstreamRejected, "bad hero sms getNumberV2 response", false)
	}
	return heroSMSProviderOrder(orderID, payload.PhoneNumber, payload, request), nil
}

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
		Price:                    core.Money{CurrencyCode: heroSMSCurrencyCode(payload.Currency), AmountDecimal: firstHeroSMSScalar(payload.ActivationCost)},
		AcquiredAt:               parseHeroSMSTime(payload.ActivationTime),
		ExpiresAt:                parseHeroSMSTime(payload.ActivationEndTime),
		CanRequestAdditionalCode: heroSMSBool(payload.CanGetAnotherSMS),
	}
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

func parseAccessNumber(result string) (orderID, rawPhone string, ok bool) {
	parts := strings.SplitN(result, ":", 3)
	if len(parts) != 3 || parts[0] != "ACCESS_NUMBER" {
		return "", "", false
	}
	return parts[1], parts[2], true
}

func heroSMSBool(raw json.RawMessage) bool {
	scalar := strings.ToLower(firstHeroSMSScalar(raw))
	return scalar == "true" || scalar == "1"
}

func heroSMSCurrencyCode(raw json.RawMessage) string {
	switch firstHeroSMSScalar(raw) {
	case "840":
		return "USD"
	default:
		return ""
	}
}

func parseHeroSMSTime(value string) time.Time {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}
	}
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05",
	}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed.UTC()
		}
	}
	return time.Time{}
}
