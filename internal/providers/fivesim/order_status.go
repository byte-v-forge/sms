package fivesim

import (
	"context"
	"net/url"

	"github.com/byte-v-forge/sms/internal/core"
)

func (c *Client) GetStatus(ctx context.Context, upstreamOrderID string) (core.ProviderCodeResult, error) {
	var payload order
	if err := c.getJSON(ctx, "/v1/user/check/"+url.PathEscape(upstreamOrderID), nil, true, &payload); err != nil {
		return core.ProviderCodeResult{}, err
	}
	return orderToCodeResult(payload), nil
}
