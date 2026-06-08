package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/eventoutbox"
)

func (s *OrderService) orderAcquiredRecords(ctx context.Context, order core.Order, reason string) ([]eventoutbox.Record, error) {
	acquired, err := s.events.OrderAcquired(ctx, order)
	if err != nil {
		return nil, err
	}
	poll, err := s.events.OrderPollRequested(ctx, order, reason)
	if err != nil {
		return nil, err
	}
	return nonEmptyRecords(acquired, poll), nil
}
