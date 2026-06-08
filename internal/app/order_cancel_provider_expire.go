package app

import (
	"context"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) expireLoadedOrder(ctx context.Context, order core.Order, now time.Time) (core.Order, error) {
	previousStatus := order.Status
	order.Status = core.StatusExpired
	order.UpdatedAt = now
	records, err := s.statusChangedRecords(ctx, order, previousStatus)
	if err != nil {
		return order, err
	}
	if err := s.updateOrder(ctx, order, records...); err != nil {
		return order, err
	}
	return order, core.NewError(core.CodeOrderExpired, "order expired", false)
}
