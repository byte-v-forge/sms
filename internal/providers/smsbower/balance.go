package smsbower

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (c *Client) GetBalance(ctx context.Context) (core.Money, error) {
	return c.api.GetBalance(ctx)
}
