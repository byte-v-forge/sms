package herosms

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/providers/handlerapi"
	"github.com/byte-v-forge/sms/internal/providers/providerhttp"
)

func (c *Client) getOpenAPIJSON(ctx context.Context, path string, params url.Values, out any) error {
	response, err := providerhttp.Do(ctx, c.httpClient, func(ctx context.Context) (*http.Request, error) {
		endpoint := c.openAPIEndpointWithPath(path)
		if len(params) > 0 {
			endpoint.RawQuery = params.Encode()
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
		if err != nil {
			return nil, core.NewError(core.CodeInternal, err.Error(), false)
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authorization", "ApiKey "+c.apiKey)
		req.Header.Set("User-Agent", c.userAgent)
		return req, nil
	}, heroSMSOpenAPIRetryPolicy())
	if err != nil {
		var smsErr *core.Error
		if errors.As(err, &smsErr) {
			return smsErr
		}
		return core.NewError(core.CodeSupplyUnavailable, err.Error(), true)
	}
	text := strings.TrimSpace(string(response.Body))
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return mapHeroSMSOpenAPIError(response.StatusCode, text)
	}
	if err := json.Unmarshal(response.Body, out); err != nil {
		return mapHeroSMSOpenAPIError(response.StatusCode, text)
	}
	return nil
}

func (c *Client) openAPIEndpointWithPath(path string) url.URL {
	endpoint := c.openAPIEndpoint
	endpoint.Path = strings.TrimRight(endpoint.Path, "/") + path
	return endpoint
}

func heroSMSOpenAPIRetryPolicy() providerhttp.RetryPolicy {
	policy := providerhttp.DefaultRetry()
	policy.MaxBodyBytes = 8 << 20
	return policy
}

func mapHeroSMSOpenAPIError(statusCode int, text string) error {
	var payload struct {
		Title   string `json:"title"`
		Details string `json:"details"`
	}
	if err := json.Unmarshal([]byte(text), &payload); err == nil && strings.TrimSpace(payload.Title) != "" {
		return handlerapi.MapTextError(text)
	}
	if text != "" {
		return handlerapi.MapTextError(text)
	}
	if statusCode >= 500 {
		return core.NewError(core.CodeSupplyUnavailable, fmt.Sprintf("hero sms openapi http status %d", statusCode), true)
	}
	return core.NewError(core.CodeUpstreamRejected, fmt.Sprintf("hero sms openapi http status %d", statusCode), false)
}
