package handlerapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/jsonx"
	"github.com/byte-v-forge/sms/internal/platform/stringx"
	"github.com/byte-v-forge/sms/internal/providers/providerhttp"
)

type HTTPDoer = providerhttp.HTTPDoer

type Client struct {
	endpoint   string
	apiKey     string
	httpClient HTTPDoer
	userAgent  string
}

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
	response, err := providerhttp.Do(ctx, c.httpClient, func(ctx context.Context) (*http.Request, error) {
		endpoint, err := url.Parse(c.endpoint)
		if err != nil {
			return nil, core.NewError(core.CodeValidationFailed, "invalid handler api endpoint", false)
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
			return nil, core.NewError(core.CodeInternal, err.Error(), false)
		}
		req.Header.Set("User-Agent", c.userAgent)
		return req, nil
	}, handlerAPIRetryPolicy(action))
	if err != nil {
		var smsErr *core.Error
		if errors.As(err, &smsErr) {
			return "", smsErr
		}
		return "", core.NewError(core.CodeSupplyUnavailable, err.Error(), true)
	}
	text := strings.TrimSpace(string(response.Body))
	if response.StatusCode < 200 || response.StatusCode > 299 {
		if text != "" {
			return "", MapTextError(text)
		}
		return "", core.NewError(core.CodeSupplyUnavailable, fmt.Sprintf("handler api http status %d", response.StatusCode), true)
	}
	return text, nil
}

func handlerAPIRetryPolicy(action string) providerhttp.RetryPolicy {
	switch action {
	case "getNumberV2", "setStatus":
		policy := providerhttp.NoRetry()
		policy.MaxBodyBytes = 1 << 20
		return policy
	default:
		policy := providerhttp.DefaultRetry()
		policy.MaxBodyBytes = 1 << 20
		return policy
	}
}

func MapTextError(text string) error {
	text = strings.TrimSpace(text)
	code, message := normalizeHandlerAPIError(text)
	switch {
	case text == "":
		return core.NewError(core.CodeUpstreamRejected, "empty upstream response", true)
	case code == "BAD_KEY":
		return core.NewError(core.CodeUpstreamRejected, "provider credential rejected", false)
	case code == "BAD_ACTION":
		return core.NewError(core.CodeUnsupportedOperation, "provider action rejected", false)
	case code == "BAD_SERVICE", code == "BAD_COUNTRY", code == "BAD_STATUS", code == "WRONG_EXCEPTION_PHONE", code == "WRONG_ACTIVATION_ID":
		return core.NewError(core.CodeValidationFailed, message, false)
	case code == "NO_ACTIVATION":
		return core.NewError(core.CodeOrderNotFound, "upstream order not found", false)
	case code == "NO_BALANCE", code == "NO_BALANCE_FORWARD":
		return core.NewError(core.CodeInsufficientBalance, "provider balance is insufficient", false)
	case code == "NO_NUMBERS", code == "NO_NUMBER", strings.Contains(text, "NO_NUMBERS"):
		return core.NewError(core.CodeNoNumberAvailable, "no upstream number available", true)
	case code == "EARLY_CANCEL_DENIED":
		return core.NewError(core.CodeCancelNotAllowed, "upstream denied early cancel", true)
	case code == "WRONG_MAX_PRICE", code == "ERROR_SQL", code == "ERROR_SQL25", code == "SERVER_ERROR":
		return core.NewError(core.CodeSupplyUnavailable, message, true)
	case code == "BANNED", code == "CHANNELS_LIMIT":
		return core.NewError(core.CodeSupplyUnavailable, message, false)
	case code == "SERVICE_NOT_AVAILABLE", code == "NOT_AVAILABLE", code == "WHATSAPP_NOT_AVAILABLE":
		return core.NewError(core.CodeSupplyUnavailable, message, true)
	default:
		return core.NewError(core.CodeUpstreamRejected, message, false)
	}
}

type handlerAPIErrorPayload struct {
	Title   string                     `json:"title"`
	Detail  string                     `json:"detail"`
	Details string                     `json:"details"`
	Info    map[string]json.RawMessage `json:"info"`
}

func normalizeHandlerAPIError(text string) (string, string) {
	text = strings.TrimSpace(text)
	if text == "" {
		return "", ""
	}
	if code, message, ok := parseHandlerAPIJSONError(text); ok {
		return code, message
	}
	code := text
	if idx := strings.Index(code, ":"); idx >= 0 {
		code = strings.TrimSpace(code[:idx])
	}
	return code, truncateHandlerAPIErrorMessage(text)
}

func parseHandlerAPIJSONError(text string) (string, string, bool) {
	var payload handlerAPIErrorPayload
	if err := json.Unmarshal([]byte(text), &payload); err != nil {
		return "", "", false
	}
	code := strings.TrimSpace(payload.Title)
	if code == "" {
		return "", "", false
	}
	parts := []string{code}
	if details := stringx.FirstNonEmpty(payload.Details, payload.Detail); details != "" {
		parts = append(parts, details)
	}
	for _, key := range []string{"min", "max"} {
		if value := jsonx.Scalar(payload.Info[key]); value != "" {
			parts = append(parts, key+"="+value)
		}
	}
	return code, truncateHandlerAPIErrorMessage(strings.Join(parts, ": ")), true
}

func truncateHandlerAPIErrorMessage(message string) string {
	return providerhttp.SanitizeErrorText(message)
}
