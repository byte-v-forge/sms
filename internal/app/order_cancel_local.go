package app

import (
	"context"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) cancelLocalOrder(ctx context.Context, order core.Order, now time.Time, clearCancelAllowed bool) (core.Order, error) {
	previousStatus := order.Status
	order.Status = core.StatusCanceled
	order.UpdatedAt = now
	order.LastError = nil
	if clearCancelAllowed {
		order.CancelAllowedAt = time.Time{}
	}
	records, err := s.statusChangedRecords(ctx, order, previousStatus)
	if err != nil {
		return order, err
	}
	if err := s.updateOrder(ctx, order, records...); err != nil {
		return core.Order{}, err
	}
	return order, nil
}
