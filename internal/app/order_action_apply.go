package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) applyAction(ctx context.Context, orderID, _ string, action core.ProviderAction, next core.OrderStatus) (core.Order, error) {
	order, err := s.loadActionOrder(ctx, orderID)
	if err != nil {
		return order, err
	}
	previousStatus := order.Status
	provider, err := s.provider(order.ProviderKey)
	if err != nil {
		return core.Order{}, err
	}
	bindOrderProviderConfig(provider, order)
	if err := provider.SetStatus(ctx, order.UpstreamOrderID, action); err != nil {
		return s.recordActionError(ctx, order, err)
	}
	return s.completeAction(ctx, order, previousStatus, action, next)
}
