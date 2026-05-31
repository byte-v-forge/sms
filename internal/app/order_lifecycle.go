package app

import (
	"context"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func orderHasCode(order core.Order) bool {
	return order.Status == core.StatusCodeReceived
}

func markCancelRequest(order core.Order, requestID string, retryAt time.Time, now time.Time) core.Order {
	order.UpdatedAt = now
	order.CancelAllowedAt = retryAt
	order.LastError = nil
	return order
}

func (s *OrderService) expireIfNeeded(ctx context.Context, order *core.Order) error {
	now := s.clock.Now()
	if !order.IsExpired(now) {
		return nil
	}
	previousStatus := order.Status
	order.Status = core.StatusExpired
	order.UpdatedAt = now
	records, err := s.statusChangedRecords(ctx, *order, previousStatus)
	if err != nil {
		return err
	}
	if err := s.updateOrder(ctx, *order, records...); err != nil {
		return err
	}
	return core.NewError(core.CodeOrderExpired, "order expired", false)
}
