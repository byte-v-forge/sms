package fivesim

import (
	"context"
	"encoding/json"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/jsonx"
)

func (c *Client) GetBalance(ctx context.Context) (core.Money, error) {
	var payload struct {
		Balance json.RawMessage `json:"balance"`
	}
	if err := c.getJSON(ctx, "/v1/user/profile", nil, true, &payload); err != nil {
		return core.Money{}, err
	}
	return core.Money{CurrencyCode: c.currencyCode, AmountDecimal: jsonx.Scalar(payload.Balance)}, nil
}
