package app

import (
	"context"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) CancelOrder(ctx context.Context, orderID, requestID string) (core.Order, error) {
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
	if orderHasCode(order) {
		return order, nil
	}
	if !order.Status.HasProviderLease() || strings.TrimSpace(order.UpstreamOrderID) == "" {
		return s.cancelLocalOrder(ctx, order, s.clock.Now(), false)
	}
	return s.queueCancelRequest(ctx, order, requestID, "api_request")
}

func (s *OrderService) RunCancelRequest(ctx context.Context, orderID, requestID string) (core.Order, error) {
	order, err := s.store.Get(ctx, orderID)
	if err != nil {
		return core.Order{}, err
	}
	return s.cancelLoadedOrder(ctx, order, requestID)
}
