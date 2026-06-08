package fivesim

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/providers/providerhttp"
)

func (c *Client) getJSON(ctx context.Context, path string, params url.Values, authenticated bool, out any) error {
	response, err := providerhttp.Do(ctx, c.httpClient, func(ctx context.Context) (*http.Request, error) {
		endpoint, err := url.Parse(c.endpoint + path)
		if err != nil {
			return nil, core.NewError(core.CodeValidationFailed, "invalid 5sim endpoint", false)
		}
		if len(params) > 0 {
			endpoint.RawQuery = params.Encode()
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
		if err != nil {
			return nil, core.NewError(core.CodeInternal, err.Error(), false)
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", c.userAgent)
		if authenticated {
			req.Header.Set("Authorization", "Bearer "+c.token)
		}
		return req, nil
	}, fiveSimRetryPolicy(path))
	if err != nil {
		var smsErr *core.Error
		if errors.As(err, &smsErr) {
			return smsErr
		}
		return core.NewError(core.CodeSupplyUnavailable, err.Error(), true)
	}
	text := strings.TrimSpace(string(response.Body))
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return mapError(response.StatusCode, text)
	}
	if err := json.Unmarshal(response.Body, out); err != nil {
		return mapError(response.StatusCode, text)
	}
	return nil
}

func fiveSimRetryPolicy(path string) providerhttp.RetryPolicy {
	var policy providerhttp.RetryPolicy
	if fiveSimRequestRetryable(path) {
		policy = providerhttp.DefaultRetry()
	} else {
		policy = providerhttp.NoRetry()
	}
	policy.MaxBodyBytes = 1 << 20
	return policy
}

func fiveSimRequestRetryable(path string) bool {
	return strings.HasPrefix(path, "/v1/guest/") || strings.HasPrefix(path, "/v1/user/check/") || path == "/v1/user/profile"
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
