package handlerapi

import (
	"context"
	"net/http"
	"net/url"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/providers/providerhttp"
)

func (c *Client) requestFactory(action string, params url.Values) providerhttp.RequestFactory {
	return func(ctx context.Context) (*http.Request, error) {
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
	}
}
