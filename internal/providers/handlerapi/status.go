package handlerapi

import (
	"context"
	"net/url"

	"github.com/byte-v-forge/sms/internal/core"
)

func (c *Client) GetStatus(ctx context.Context, upstreamOrderID string, parse StatusParser) (core.ProviderCodeResult, error) {
	params := url.Values{}
	params.Set("id", upstreamOrderID)
	result, err := c.Do(ctx, "getStatus", params)
	if err != nil {
		return core.ProviderCodeResult{}, err
	}
	return parse(result)
}
