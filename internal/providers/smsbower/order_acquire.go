package smsbower

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/providers/handlerapi"
)

func (c *Client) AcquireNumber(ctx context.Context, request core.ProviderAcquireRequest) (core.ProviderOrder, error) {
	result, err := c.api.GetNumberV2(ctx, request, handlerapi.GetNumberV2Config{
		ProviderName:       "smsbower",
		ProviderIDParam:    "providerIds",
		ProviderIDRequired: true,
	})
	if err != nil {
		return core.ProviderOrder{}, err
	}
	order, err := c.parseGetNumberV2(result, request)
	if err == nil {
		return order, nil
	}
	if isProviderTextError(result) {
		return core.ProviderOrder{}, handlerapi.MapTextError(result)
	}
	return core.ProviderOrder{}, err
}
