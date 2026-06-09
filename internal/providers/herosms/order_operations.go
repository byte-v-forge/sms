package herosms

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (c *Client) AcquireNumber(ctx context.Context, request core.ProviderAcquireRequest) (core.ProviderOrder, error) {
	return c.buyActivation(ctx, request)
}
