package handlerapi

import (
	"context"
	"net/url"

	"github.com/byte-v-forge/sms/internal/core"
)

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

func (c *Client) SetActivationStatus(ctx context.Context, upstreamOrderID string, action core.ProviderAction, providerName string) error {
	return c.SetStatus(ctx, upstreamOrderID, action, func(action core.ProviderAction) (string, string, error) {
		return ActivationStatusForAction(providerName, action)
	})
}
