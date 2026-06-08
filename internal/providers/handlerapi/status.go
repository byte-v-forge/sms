package handlerapi

import (
	"context"
	"net/url"
	"strings"

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

func ActivationStatusForAction(providerName string, action core.ProviderAction) (status string, expected string, err error) {
	switch action {
	case core.ActionMarkMessageSent:
		return "1", "ACCESS_READY", nil
	case core.ActionRequestAdditional:
		return "3", "ACCESS_RETRY_GET", nil
	case core.ActionCompleteOrder:
		return "6", "ACCESS_ACTIVATION", nil
	case core.ActionCancelOrder:
		return "8", "ACCESS_CANCEL", nil
	default:
		providerName = strings.TrimSpace(providerName)
		if providerName == "" {
			providerName = "sms provider"
		}
		return "", "", core.NewError(core.CodeUnsupportedOperation, "unsupported "+providerName+" status action", false)
	}
}
