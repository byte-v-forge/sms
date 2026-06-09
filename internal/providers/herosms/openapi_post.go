package herosms

import (
	"context"
	"errors"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/providers/providerhttp"
)

func (c *Client) postOpenAPIJSON(ctx context.Context, path string, body any) ([]byte, error) {
	response, err := providerhttp.Do(ctx, c.httpClient, c.openAPIPostRequestFactory(path, body), heroSMSOpenAPINonRetryPolicy())
	if err != nil {
		var smsErr *core.Error
		if errors.As(err, &smsErr) {
			return nil, smsErr
		}
		return nil, core.NewError(core.CodeSupplyUnavailable, err.Error(), true)
	}
	text := strings.TrimSpace(string(response.Body))
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, mapHeroSMSOpenAPIError(response.StatusCode, text)
	}
	return response.Body, nil
}
