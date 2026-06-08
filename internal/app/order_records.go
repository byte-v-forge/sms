package app

import (
	"context"
	"strings"

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

func (s *OrderService) statusAndCodeRecords(ctx context.Context, order core.Order, previous core.OrderStatus, code core.SMSCode) ([]eventoutbox.Record, error) {
	records, err := s.statusChangedRecords(ctx, order, previous)
	if err != nil {
		return nil, err
	}
	codeRecord, err := s.events.CodeReceived(ctx, order, code)
	if err != nil {
		return nil, err
	}
	return nonEmptyRecords(append(records, codeRecord)...), nil
}

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

func nonEmptyRecords(records ...eventoutbox.Record) []eventoutbox.Record {
	out := make([]eventoutbox.Record, 0, len(records))
	for _, record := range records {
		if strings.TrimSpace(record.EventID) == "" {
			continue
		}
		out = append(out, record)
	}
	return out
}
