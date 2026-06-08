package handlerapi

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/providers/providerhttp"
)

func (c *Client) Do(ctx context.Context, action string, params url.Values) (string, error) {
	response, err := providerhttp.Do(ctx, c.httpClient, c.requestFactory(action, params), handlerAPIRetryPolicy(action))
	if err != nil {
		var smsErr *core.Error
		if errors.As(err, &smsErr) {
			return "", smsErr
		}
		return "", core.NewError(core.CodeSupplyUnavailable, err.Error(), true)
	}
	text := strings.TrimSpace(string(response.Body))
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return "", handlerAPIHTTPError(response.StatusCode, text)
	}
	return text, nil
}
