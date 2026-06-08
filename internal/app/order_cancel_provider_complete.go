package app

import (
	"context"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) completeProviderCancel(ctx context.Context, order core.Order, now time.Time) (core.Order, error) {
	previousStatus := order.Status
	order.Status = core.StatusCanceled
	order.UpdatedAt = now
	order.LastError = nil
	order.CancelAllowedAt = time.Time{}
	records, err := s.statusChangedRecords(ctx, order, previousStatus)
	if err != nil {
		return core.Order{}, err
	}
	if err := s.updateOrder(ctx, order, records...); err != nil {
		return core.Order{}, err
	}
	return order, nil
}
