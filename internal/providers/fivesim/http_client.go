package fivesim

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/providers/providerhttp"
)

func (c *Client) getJSON(ctx context.Context, path string, params url.Values, authenticated bool, out any) error {
	response, err := providerhttp.Do(ctx, c.httpClient, c.requestFactory(path, params, authenticated), fiveSimRetryPolicy(path))
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
