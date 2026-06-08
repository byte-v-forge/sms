package handlerapi

import (
	"context"
	"net/url"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

type StatusParser func(string) (core.ProviderCodeResult, error)

type StatusActionMapper func(core.ProviderAction) (status string, expected string, err error)

func (c *Client) GetStatus(ctx context.Context, upstreamOrderID string, parse StatusParser) (core.ProviderCodeResult, error) {
	params := url.Values{}
	params.Set("id", upstreamOrderID)
	result, err := c.Do(ctx, "getStatus", params)
	if err != nil {
		return core.ProviderCodeResult{}, err
	}
	return parse(result)
}

func (c *Client) SetStatus(ctx context.Context, upstreamOrderID string, action core.ProviderAction, mapAction StatusActionMapper) error {
	status, expected, err := mapAction(action)
	if err != nil {
		return err
	}
	params := url.Values{}
	params.Set("id", upstreamOrderID)
	params.Set("status", status)
	result, err := c.Do(ctx, "setStatus", params)
	if err != nil {
		return err
	}
	if result != expected {
		return MapTextError(result)
	}
	return nil
}

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
