package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/eventoutbox"
)

func (s *OrderService) statusChangedRecords(ctx context.Context, order core.Order, previous core.OrderStatus) ([]eventoutbox.Record, error) {
	if previous == order.Status {
		return nil, nil
	}
	record, err := s.events.OrderStatusChanged(ctx, order, previous)
	if err != nil {
		return nil, err
	}
	return nonEmptyRecords(record), nil
}
