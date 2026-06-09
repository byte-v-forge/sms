package fivesim

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (c *Client) AcquireNumber(ctx context.Context, request core.ProviderAcquireRequest) (core.ProviderOrder, error) {
	path, err := acquireNumberPath(request.Route)
	if err != nil {
		return core.ProviderOrder{}, err
	}
	var payload order
	if err := c.getJSON(ctx, path, nil, true, &payload); err != nil {
		return core.ProviderOrder{}, err
	}
	return providerOrderFromPayload(payload, request, c.currencyCode)
}
