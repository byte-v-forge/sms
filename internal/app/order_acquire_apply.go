package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) applyProviderOrder(ctx context.Context, order core.Order, route core.Route, provider core.Provider, providerOrder core.ProviderOrder) (core.Order, error) {
	now := s.clock.Now()
	previousStatus := order.Status
	order = providerOrderApplied(order, route, provider, providerOrder, now)
	records, err := s.orderAcquiredRecords(ctx, order, "order_acquired")
	if err != nil {
		return core.Order{}, err
	}
	statusRecords, err := s.statusChangedRecords(ctx, order, previousStatus)
	if err != nil {
		return core.Order{}, err
	}
	records = append(records, statusRecords...)
	if err := s.updateOrder(ctx, order, records...); err != nil {
		return core.Order{}, err
	}
	return order, nil
}
