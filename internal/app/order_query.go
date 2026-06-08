package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) GetOrder(ctx context.Context, orderID string) (core.Order, error) {
	order, err := s.store.Get(ctx, orderID)
	if err != nil {
		return core.Order{}, err
	}
	if err := s.expireIfNeeded(ctx, &order); err != nil {
		return core.Order{}, err
	}
	if s.execution.SyncProviderStateOnRead() && order.Status.HasProviderLease() && order.Status != core.StatusCodeReceived && !order.Status.IsFinal() {
		synced, _, err := s.CheckCode(ctx, orderID)
		if synced.ID == "" {
			synced = order
		}
		if err == nil {
			return synced, nil
		}
		return synced, err
	}
	return order, nil
}
