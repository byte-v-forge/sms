package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) recordAcquireFailure(ctx context.Context, order core.Order, smsErr *core.Error) (core.Order, error) {
	now := s.clock.Now()
	order.LastError = smsErr
	order.UpdatedAt = now
	if smsErr.Retryable {
		if err := s.updateOrder(ctx, order); err != nil {
			return core.Order{}, err
		}
		return order, smsErr
	}
	previousStatus := order.Status
	order.Status = core.StatusFailed
	records, err := s.statusChangedRecords(ctx, order, previousStatus)
	if err != nil {
		return core.Order{}, err
	}
	if err := s.updateOrder(ctx, order, records...); err != nil {
		return core.Order{}, err
	}
	return order, smsErr
}
