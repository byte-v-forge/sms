package handlerapi

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	endpoint   string
	apiKey     string
	httpClient HTTPDoer
	userAgent  string
}

const maxResponseBytes = 16 << 20

func New(endpoint, apiKey string, httpClient HTTPDoer) (*Client, error) {
	if strings.TrimSpace(endpoint) == "" {
		return nil, core.NewError(core.CodeValidationFailed, "handler api endpoint is required", false)
	}
	if strings.TrimSpace(apiKey) == "" {
		return nil, core.NewError(core.CodeValidationFailed, "handler api key is required", false)
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}
	return &Client{
		endpoint:   endpoint,
		apiKey:     apiKey,
		httpClient: httpClient,
		userAgent:  "sms/1.0",
	}, nil
}

func (c *Client) Do(ctx context.Context, action string, params url.Values) (string, error) {
	endpoint, err := url.Parse(c.endpoint)
	if err != nil {
		return "", core.NewError(core.CodeValidationFailed, "invalid handler api endpoint", false)
	}
	query := endpoint.Query()
	for key, values := range params {
		for _, value := range values {
			if value != "" {
				query.Add(key, value)
			}
		}
	}
	query.Set("api_key", c.apiKey)
	query.Set("action", action)
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return "", core.NewError(core.CodeInternal, err.Error(), false)
	}
	req.Header.Set("User-Agent", c.userAgent)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", core.NewError(core.CodeSupplyUnavailable, err.Error(), true)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes))
	if err != nil {
		return "", core.NewError(core.CodeSupplyUnavailable, err.Error(), true)
	}
	text := strings.TrimSpace(string(body))
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", core.NewError(core.CodeSupplyUnavailable, fmt.Sprintf("handler api http status %d", resp.StatusCode), true)
	}
	return text, nil
}

func MapTextError(text string) error {
	text = strings.TrimSpace(text)
	code := text
	if idx := strings.Index(code, ":"); idx >= 0 {
		code = strings.TrimSpace(code[:idx])
	}
	switch {
	case text == "":
		return core.NewError(core.CodeUpstreamRejected, "empty upstream response", true)
	case code == "BAD_KEY":
		return core.NewError(core.CodeUpstreamRejected, "provider credential rejected", false)
	case code == "BAD_ACTION":
		return core.NewError(core.CodeUnsupportedOperation, "provider action rejected", false)
	case code == "BAD_SERVICE", code == "BAD_COUNTRY", code == "BAD_STATUS", code == "WRONG_EXCEPTION_PHONE", code == "WRONG_ACTIVATION_ID":
		return core.NewError(core.CodeValidationFailed, text, false)
	case code == "NO_ACTIVATION":
		return core.NewError(core.CodeActivationNotFound, "upstream activation not found", false)
	case code == "NO_BALANCE", code == "NO_BALANCE_FORWARD":
		return core.NewError(core.CodeInsufficientBalance, "provider balance is insufficient", false)
	case code == "WRONG_MAX_PRICE", code == "BAD_MAX_PRICE":
		return core.NewError(core.CodePriceLimitExceeded, text, false)
	case code == "NO_NUMBERS", code == "NO_NUMBER", strings.Contains(text, "NO_NUMBERS"):
		return core.NewError(core.CodeNoNumberAvailable, "no upstream number available", true)
	case code == "EARLY_CANCEL_DENIED":
		return core.NewError(core.CodeCancelNotAllowed, "upstream denied early cancel", true)
	case code == "ERROR_SQL", code == "ERROR_SQL25", code == "SERVER_ERROR":
		return core.NewError(core.CodeSupplyUnavailable, text, true)
	case code == "BANNED", code == "CHANNELS_LIMIT":
		return core.NewError(core.CodeSupplyUnavailable, text, false)
	case code == "SERVICE_NOT_AVAILABLE", code == "NOT_AVAILABLE", code == "WHATSAPP_NOT_AVAILABLE":
		return core.NewError(core.CodeSupplyUnavailable, text, true)
	default:
		return core.NewError(core.CodeUpstreamRejected, text, false)
	}
}
