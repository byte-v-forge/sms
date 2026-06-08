package herosms

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/providers/providerhttp"
)

func (c *Client) getOpenAPIJSON(ctx context.Context, path string, params url.Values, out any) error {
	response, err := providerhttp.Do(ctx, c.httpClient, c.openAPIRequestFactory(path, params), heroSMSOpenAPIRetryPolicy())
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
