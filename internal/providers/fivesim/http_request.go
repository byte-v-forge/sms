package fivesim

import (
	"context"
	"net/http"
	"net/url"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/providers/providerhttp"
)

func (c *Client) requestFactory(path string, params url.Values, authenticated bool) providerhttp.RequestFactory {
	return func(ctx context.Context) (*http.Request, error) {
		endpoint := c.endpointWithPath(path)
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
	}
}
