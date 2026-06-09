package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) applySyncedProviderTerminalStatus(
	ctx context.Context,
	order core.Order,
	previousStatus core.OrderStatus,
	status core.OrderStatus,
) (core.Order, error) {
	order.Status = status
	records, err := s.statusChangedRecords(ctx, order, previousStatus)
	if err != nil {
		return order, err
	}
	if err := s.updateOrder(ctx, order, records...); err != nil {
		return core.Order{}, err
	}
	return order, nil
}
