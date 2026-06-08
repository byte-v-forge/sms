package handlerapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/providers/providerhttp"
)

type HTTPDoer = providerhttp.HTTPDoer

type Client struct {
	endpoint   url.URL
	apiKey     string
	httpClient HTTPDoer
	userAgent  string
}

func New(rawEndpoint, apiKey string, httpClient HTTPDoer) (*Client, error) {
	rawEndpoint = strings.TrimSpace(rawEndpoint)
	if rawEndpoint == "" {
		return nil, core.NewError(core.CodeValidationFailed, "handler api endpoint is required", false)
	}
	endpoint, err := url.Parse(rawEndpoint)
	if err != nil {
		return nil, core.NewError(core.CodeValidationFailed, "invalid handler api endpoint", false)
	}
	if strings.TrimSpace(apiKey) == "" {
		return nil, core.NewError(core.CodeValidationFailed, "handler api key is required", false)
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}
	return &Client{
		endpoint:   *endpoint,
		apiKey:     apiKey,
		httpClient: httpClient,
		userAgent:  "sms/1.0",
	}, nil
}

func (c *Client) Do(ctx context.Context, action string, params url.Values) (string, error) {
	response, err := providerhttp.Do(ctx, c.httpClient, func(ctx context.Context) (*http.Request, error) {
		endpoint := c.endpoint
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
