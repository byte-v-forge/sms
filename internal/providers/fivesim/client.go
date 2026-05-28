package fivesim

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/byte-v-forge/common-lib/httpx"
	"github.com/byte-v-forge/common-lib/stringx"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/providers/handlerapi"
	"github.com/byte-v-forge/sms/internal/providers/phone"
)

const (
	DefaultEndpoint = "https://5sim.net"
	ProviderKey     = "5sim"
)

type Config struct {
	Endpoint         string
	Token            string
	CurrencyCode     string
	Ref              string
	Reuse            bool
	Voice            bool
	Forwarding       bool
	ForwardingNumber string
}

type Client struct {
	endpoint         string
	token            string
	currencyCode     string
	ref              string
	reuse            bool
	voice            bool
	forwarding       bool
	forwardingNumber string
	httpClient       handlerapi.HTTPDoer
	userAgent        string
	policy           core.ProviderPolicy
}

type order struct {
	ID               json.RawMessage `json:"id"`
	CreatedAt        string          `json:"created_at"`
	Phone            string          `json:"phone"`
	Operator         string          `json:"operator"`
	Product          string          `json:"product"`
	Price            json.RawMessage `json:"price"`
	Status           string          `json:"status"`
	Expires          string          `json:"expires"`
	SMS              []sms           `json:"sms"`
	Forwarding       bool            `json:"forwarding"`
	ForwardingNumber string          `json:"forwarding_number"`
	Country          string          `json:"country"`
}

type sms struct {
	ID        json.RawMessage `json:"id"`
	CreatedAt string          `json:"created_at"`
	Date      string          `json:"date"`
	Sender    string          `json:"sender"`
	Text      string          `json:"text"`
	Code      string          `json:"code"`
}

func New(config Config, httpClient handlerapi.HTTPDoer) (*Client, error) {
	endpoint := strings.TrimRight(strings.TrimSpace(config.Endpoint), "/")
	if endpoint == "" {
		endpoint = DefaultEndpoint
	}
	if strings.TrimSpace(config.Token) == "" {
		return nil, core.NewError(core.CodeValidationFailed, "5sim token is required", false)
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}
	return &Client{
		endpoint:         endpoint,
		token:            strings.TrimSpace(config.Token),
		currencyCode:     strings.TrimSpace(config.CurrencyCode),
		ref:              strings.TrimSpace(config.Ref),
		reuse:            config.Reuse,
		voice:            config.Voice,
		forwarding:       config.Forwarding,
		forwardingNumber: strings.TrimSpace(config.ForwardingNumber),
		httpClient:       httpClient,
		userAgent:        "sms/1.0",
		policy: core.ProviderPolicy{
			OrderTTL:     20 * time.Minute,
			PollInterval: 5 * time.Second,
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
	case core.ActionMarkMessageSent, core.ActionRequestAdditional:
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

func (c *Client) getJSON(ctx context.Context, path string, params url.Values, authenticated bool, out any) error {
	endpoint, err := url.Parse(c.endpoint + path)
	if err != nil {
		return core.NewError(core.CodeValidationFailed, "invalid 5sim endpoint", false)
	}
	if len(params) > 0 {
		endpoint.RawQuery = params.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return core.NewError(core.CodeInternal, err.Error(), false)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	if authenticated {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return core.NewError(core.CodeSupplyUnavailable, err.Error(), true)
	}
	defer resp.Body.Close()
	body, err := httpx.ReadLimited(resp.Body, 1<<20)
	if err != nil {
		return core.NewError(core.CodeSupplyUnavailable, err.Error(), true)
	}
	text := strings.TrimSpace(string(body))
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return mapError(resp.StatusCode, text)
	}
	if err := json.Unmarshal(body, out); err != nil {
		return mapError(resp.StatusCode, text)
	}
	return nil
}

func orderToCodeResult(payload order) core.ProviderCodeResult {
	latestSMS := latestSMS(payload.SMS)
	if latestSMS.Code != "" || latestSMS.Text != "" {
		receivedAt := parseTime(stringx.FirstNonEmpty(latestSMS.Date, latestSMS.CreatedAt))
		return core.ProviderCodeResult{
			Status:      core.StatusCodeReceived,
			Code:        latestSMS.Code,
			MessageText: latestSMS.Text,
			ReceivedAt:  receivedAt,
		}
	}
	switch strings.ToUpper(strings.TrimSpace(payload.Status)) {
	case "CANCELED":
		return core.ProviderCodeResult{Status: core.StatusCanceled}
	case "TIMEOUT":
		return core.ProviderCodeResult{Status: core.StatusExpired}
	case "FINISHED":
		return core.ProviderCodeResult{Status: core.StatusCompleted}
	case "BANNED":
		return core.ProviderCodeResult{Status: core.StatusFailed}
	default:
		return core.ProviderCodeResult{Status: core.StatusPendingCode}
	}
}

func latestSMS(messages []sms) sms {
	if len(messages) == 0 {
		return sms{}
	}
	latest := messages[0]
	latestAt := parseTime(stringx.FirstNonEmpty(latest.Date, latest.CreatedAt))
	for _, message := range messages[1:] {
		messageAt := parseTime(stringx.FirstNonEmpty(message.Date, message.CreatedAt))
		if latestAt.IsZero() || messageAt.After(latestAt) {
			latest = message
			latestAt = messageAt
		}
	}
	return latest
}

func mapError(statusCode int, text string) error {
	normalized := strings.ToLower(strings.TrimSpace(text))
	switch {
	case statusCode == http.StatusUnauthorized:
		return core.NewError(core.CodeUpstreamRejected, "5sim credential rejected", false)
	case strings.Contains(normalized, "order not found"), strings.Contains(normalized, "record not found"):
		return core.NewError(core.CodeOrderNotFound, text, false)
	case strings.Contains(normalized, "no free phones"):
		return core.NewError(core.CodeNoNumberAvailable, text, true)
	case strings.Contains(normalized, "not enough user balance"), strings.Contains(normalized, "insufficient"):
		return core.NewError(core.CodeInsufficientBalance, text, false)
	case strings.Contains(normalized, "order expired"):
		return core.NewError(core.CodeOrderExpired, text, false)
	case strings.Contains(normalized, "order has sms"):
		return core.NewError(core.CodeCancelNotAllowed, text, false)
	case strings.Contains(normalized, "bad country"), strings.Contains(normalized, "bad operator"), strings.Contains(normalized, "bad product"),
		strings.Contains(normalized, "select country"), strings.Contains(normalized, "select operator"), strings.Contains(normalized, "select product"),
		strings.Contains(normalized, "product is empty"):
		return core.NewError(core.CodeValidationFailed, text, false)
	case statusCode >= 500:
		return core.NewError(core.CodeSupplyUnavailable, text, true)
	case text == "":
		return core.NewError(core.CodeUpstreamRejected, fmt.Sprintf("5sim http status %d", statusCode), statusCode >= 500)
	default:
		return core.NewError(core.CodeUpstreamRejected, text, false)
	}
}

func parseTime(value string) time.Time {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}
	}
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05.999999Z",
		"2006-01-02 15:04:05",
	}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed.UTC()
		}
	}
	return time.Time{}
}

func rawJSONScalar(raw json.RawMessage) string {
	if len(raw) == 0 || string(raw) == "null" {
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

func setBool(values url.Values, key string, value bool) {
	if value {
		values.Set(key, "1")
	}
}
