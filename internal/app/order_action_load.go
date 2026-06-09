package app

import (
	"context"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) loadActionOrder(ctx context.Context, orderID string) (core.Order, error) {
	order, err := s.store.Get(ctx, orderID)
	if err != nil {
		return core.Order{}, err
	}
	if err := s.expireIfNeeded(ctx, &order); err != nil {
		return order, err
	}
	if order.Status.IsFinal() {
		return order, core.NewError(core.CodeOrderAlreadyFinalized, "order already finalized", false)
	}
	if !order.Status.HasProviderLease() || strings.TrimSpace(order.UpstreamOrderID) == "" {
		return order, core.NewError(core.CodeOrderNotFound, "order has no upstream provider lease", true)
	}
	return order, nil
}
