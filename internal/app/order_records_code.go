package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/eventoutbox"
)

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
