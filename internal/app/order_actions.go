package app

import (
	"context"

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
