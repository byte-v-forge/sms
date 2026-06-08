package handlerapi

import (
	"context"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func (c *Client) GetBalance(ctx context.Context) (core.Money, error) {
	result, err := c.Do(ctx, "getBalance", nil)
	if err != nil {
		return core.Money{}, err
	}
	const prefix = "ACCESS_BALANCE:"
	if !strings.HasPrefix(result, prefix) {
		return core.Money{}, MapTextError(result)
	}
	return core.Money{AmountDecimal: strings.TrimPrefix(result, prefix)}, nil
}
