package app

import (
	"context"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) queueCancelRequest(ctx context.Context, order core.Order, requestID string, reason string) (core.Order, error) {
	if requestID = strings.TrimSpace(requestID); requestID == "" {
		requestID = s.ids.NewID("req_")
	}
	retryAt := s.cancelReadyAt(ctx, order)
	order = markCancelRequest(order, requestID, retryAt, s.clock.Now())
	record, err := s.events.OrderCancelRequested(ctx, order, requestID, reason)
	if err != nil {
		return order, err
	}
	if err := s.updateOrder(ctx, order, record); err != nil {
		return core.Order{}, err
	}
	return s.execution.AfterCancelQueued(ctx, order, requestID)
}
