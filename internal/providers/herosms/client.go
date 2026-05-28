package herosms

import (
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/providers/handlerapi"
	"github.com/byte-v-forge/sms/internal/providers/phone"
)

const (
	DefaultEndpoint = "https://hero-sms.com/stubs/handler_api.php"
	ProviderKey     = "herosms"
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
			OrderTTL:           20 * time.Minute,
			PollInterval:       5 * time.Second,
			CancelAllowedAfter: 2 * time.Minute,
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

	result, err := c.api.Do(ctx, "getNumber", params)
	if err != nil {
		return core.ProviderOrder{}, err
	}
	orderID, rawPhone, ok := parseAccessNumber(result)
	if !ok {
		return core.ProviderOrder{}, handlerapi.MapTextError(result)
	}
	e164, national := phone.Normalize(rawPhone, request.Target.CountryISO2, request.Target.CountryCallingCode)
	return core.ProviderOrder{
		UpstreamOrderID: orderID,
		PhoneNumber: core.PhoneNumber{
			E164:               e164,
			NationalNumber:     national,
			CountryISO2:        request.Target.CountryISO2,
			CountryCallingCode: request.Target.CountryCallingCode,
		},
		AcquiredAt: time.Now().UTC(),
	}, nil
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

func parseAccessNumber(result string) (orderID, rawPhone string, ok bool) {
	parts := strings.SplitN(result, ":", 3)
	if len(parts) != 3 || parts[0] != "ACCESS_NUMBER" {
		return "", "", false
	}
	return parts[1], parts[2], true
}

func parseStatus(result string) (core.ProviderCodeResult, error) {
	switch {
	case strings.HasPrefix(result, "STATUS_OK:"):
		return core.ProviderCodeResult{
			Status:     core.StatusCodeReceived,
			Code:       strings.TrimSpace(strings.TrimPrefix(result, "STATUS_OK:")),
			ReceivedAt: time.Now().UTC(),
		}, nil
	case result == "STATUS_WAIT_CODE":
		return core.ProviderCodeResult{Status: core.StatusPendingCode}, nil
	case strings.HasPrefix(result, "STATUS_WAIT_RETRY"):
		return core.ProviderCodeResult{Status: core.StatusAdditionalCodeRequested}, nil
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
		return "", "", core.NewError(core.CodeUnsupportedOperation, "unsupported hero sms status action", false)
	}
}
