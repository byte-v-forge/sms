package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) providerStatus(ctx context.Context, order core.Order, provider core.Provider) (core.ProviderCodeResult, error) {
	result, err := provider.GetStatus(ctx, order.UpstreamOrderID)
	if err != nil {
		order.LastError = asCoreError(err)
		order.UpdatedAt = s.clock.Now()
		_ = s.updateOrder(ctx, order)
		return core.ProviderCodeResult{}, err
	}
	return result, nil
}
