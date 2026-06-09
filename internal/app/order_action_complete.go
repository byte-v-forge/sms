package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) completeAction(ctx context.Context, order core.Order, previousStatus core.OrderStatus, action core.ProviderAction, next core.OrderStatus) (core.Order, error) {
	order.Status = next
	order.UpdatedAt = s.clock.Now()
	records, err := s.actionRecords(ctx, order, previousStatus, action)
	if err != nil {
		return core.Order{}, err
	}
	if err := s.updateOrder(ctx, order, records...); err != nil {
		return core.Order{}, err
	}
	return order, nil
}

func (s *OrderService) recordActionError(ctx context.Context, order core.Order, err error) (core.Order, error) {
	order.LastError = asCoreError(err)
	order.UpdatedAt = s.clock.Now()
	_ = s.updateOrder(ctx, order)
	return order, err
}
