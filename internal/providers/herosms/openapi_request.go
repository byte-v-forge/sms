package herosms

import (
	"context"
	"net/http"
	"net/url"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/providers/providerhttp"
)

func (c *Client) openAPIRequestFactory(path string, params url.Values) providerhttp.RequestFactory {
	return func(ctx context.Context) (*http.Request, error) {
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
	}
}
