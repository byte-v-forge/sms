package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/eventoutbox"
)

func (s *OrderService) actionRecords(ctx context.Context, order core.Order, previous core.OrderStatus, action core.ProviderAction) ([]eventoutbox.Record, error) {
	records, err := s.statusChangedRecords(ctx, order, previous)
	if err != nil {
		return nil, err
	}
	if action != core.ActionMarkMessageSent && action != core.ActionRequestAdditional {
		return records, nil
	}
	poll, err := s.events.OrderPollRequested(ctx, order, string(action))
	if err != nil {
		return nil, err
	}
	return nonEmptyRecords(append(records, poll)...), nil
}
