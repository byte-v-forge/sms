package smsbower

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (c *Client) GetStatus(ctx context.Context, upstreamOrderID string) (core.ProviderCodeResult, error) {
	return c.api.GetStatus(ctx, upstreamOrderID, parseStatus)
}

func (c *Client) SetStatus(ctx context.Context, upstreamOrderID string, action core.ProviderAction) error {
	return c.api.SetActivationStatus(ctx, upstreamOrderID, action, "smsbower")
}
