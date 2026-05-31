package app

import (
	"context"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) MarkMessageSent(ctx context.Context, orderID, requestID string) (core.Order, error) {
	return s.applyAction(ctx, orderID, requestID, core.ActionMarkMessageSent, core.StatusMessageSent)
}

func (s *OrderService) RequestAdditionalCode(ctx context.Context, orderID, requestID string) (core.Order, error) {
	return s.applyAction(ctx, orderID, requestID, core.ActionRequestAdditional, core.StatusAdditionalCodeRequested)
}

func (s *OrderService) CompleteOrder(ctx context.Context, orderID, requestID string) (core.Order, error) {
	return s.applyAction(ctx, orderID, requestID, core.ActionCompleteOrder, core.StatusCompleted)
}

func (s *OrderService) applyAction(ctx context.Context, orderID, _ string, action core.ProviderAction, next core.OrderStatus) (core.Order, error) {
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
	previousStatus := order.Status
	provider, err := s.provider(order.ProviderKey)
	if err != nil {
		return core.Order{}, err
	}
	bindOrderProviderConfig(provider, order)
	if err := provider.SetStatus(ctx, order.UpstreamOrderID, action); err != nil {
		order.LastError = asCoreError(err)
		order.UpdatedAt = s.clock.Now()
		_ = s.updateOrder(ctx, order)
		return order, err
	}
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
