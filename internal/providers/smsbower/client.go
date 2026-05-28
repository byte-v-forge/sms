package smsbower

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/providers/handlerapi"
	"github.com/byte-v-forge/sms/internal/providers/phone"
)

const (
	DefaultEndpoint = "https://smsbower.page/stubs/handler_api.php"
	ProviderKey     = "smsbower"
)

type Config struct {
	Endpoint string
	APIKey   string
}

type Client struct {
	api    *handlerapi.Client
	policy core.ProviderPolicy
}

func New(config Config, httpClient handlerapi.HTTPDoer) (*Client, error) {
	endpoint := config.Endpoint
	if endpoint == "" {
		endpoint = DefaultEndpoint
	}
	api, err := handlerapi.New(endpoint, config.APIKey, httpClient)
	if err != nil {
		return nil, err
	}
	return &Client{
		api: api,
		policy: core.ProviderPolicy{
			OrderTTL:              25 * time.Minute,
			PollInterval:          5 * time.Second,
			EarlyCancelRetryAfter: 2 * time.Minute,
		},
	}, nil
}

func (c *Client) Key() string {
	return ProviderKey
}

func (c *Client) Policy() core.ProviderPolicy {
	return c.policy
}

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

func matchService(candidate string, applications []ApplicationOffer) string {
	normalized := normalizeApplicationAlias(candidate)
	for _, app := range applications {
		if normalizeApplicationAlias(app.UpstreamServiceKey) == normalized || normalizeApplicationAlias(app.ApplicationKey) == normalized {
			return app.UpstreamServiceKey
		}
	}
	for _, app := range applications {
		display := normalizeApplicationAlias(app.DisplayName)
		if display != "" && (display == normalized || strings.Contains(display, normalized)) {
			return app.UpstreamServiceKey
		}
	}
	return ""
}

func normalizeApplicationAlias(value string) string {
	var out strings.Builder
	for _, r := range strings.ToLower(strings.TrimSpace(value)) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			out.WriteRune(r)
		}
	}
	return out.String()
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

func parseStatus(result string) (core.ProviderCodeResult, error) {
	switch {
	case strings.HasPrefix(result, "STATUS_OK:"):
		return core.ProviderCodeResult{
			Status:     core.StatusCodeReceived,
			Code:       strings.Trim(strings.TrimSpace(strings.TrimPrefix(result, "STATUS_OK:")), "'\""),
			ReceivedAt: time.Now().UTC(),
		}, nil
	case result == "STATUS_WAIT_CODE":
		return core.ProviderCodeResult{Status: core.StatusPendingCode}, nil
	case strings.HasPrefix(result, "STATUS_WAIT_RETRY:"):
		return core.ProviderCodeResult{
			Status: core.StatusAdditionalCodeRequested,
			Code:   strings.TrimSpace(strings.TrimPrefix(result, "STATUS_WAIT_RETRY:")),
		}, nil
	case result == "STATUS_CANCEL":
		return core.ProviderCodeResult{Status: core.StatusCanceled}, nil
	default:
		return core.ProviderCodeResult{}, handlerapi.MapTextError(result)
	}
}

func statusForAction(action core.ProviderAction) (status string, expected string, err error) {
	switch action {
	case core.ActionMarkMessageSent:
		return "1", "ACCESS_READY", nil
	case core.ActionRequestAdditional:
		return "3", "ACCESS_RETRY_GET", nil
	case core.ActionCompleteOrder:
		return "6", "ACCESS_ACTIVATION", nil
	case core.ActionCancelOrder:
		return "8", "ACCESS_CANCEL", nil
	default:
		return "", "", core.NewError(core.CodeUnsupportedOperation, "unsupported smsbower status action", false)
	}
}

func isProviderTextError(result string) bool {
	return !strings.HasPrefix(strings.TrimSpace(result), "{")
}

func rawJSONScalar(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var text string
	if err := json.Unmarshal(raw, &text); err == nil {
		return strings.TrimSpace(text)
	}
	var number json.Number
	if err := json.Unmarshal(raw, &number); err == nil {
		return number.String()
	}
	var floatValue float64
	if err := json.Unmarshal(raw, &floatValue); err == nil {
		return strconv.FormatFloat(floatValue, 'f', -1, 64)
	}
	return strings.Trim(string(raw), "\"")
}

func parseOrderTime(raw json.RawMessage) time.Time {
	return parseOrderTimeText(rawJSONScalar(raw))
}

func parseOrderTimeText(value string) time.Time {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Now().UTC()
	}
	if unix, err := strconv.ParseInt(value, 10, 64); err == nil {
		if unix > 1_000_000_000_000 {
			return time.UnixMilli(unix).UTC()
		}
		return time.Unix(unix, 0).UTC()
	}
	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
	}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed.UTC()
		}
	}
	return time.Now().UTC()
}

func decodeJSONObject(result string, out any) error {
	if err := json.Unmarshal([]byte(result), out); err != nil {
		return core.NewError(core.CodeUpstreamRejected, fmt.Sprintf("bad json response: %v", err), false)
	}
	return nil
}
