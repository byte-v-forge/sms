package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) syncTerminalOrChargedProviderState(ctx context.Context, order core.Order, provider core.Provider) (core.Order, bool, error) {
	result, err := provider.GetStatus(ctx, order.UpstreamOrderID)
	if err != nil {
		return order, false, err
	}
	previousStatus := order.Status
	order = prepareSyncedProviderOrder(order, s.clock.Now())
	synced, handled, err := s.applySyncedProviderResult(ctx, order, previousStatus, result)
	return synced, handled, err
}
