package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) applyReceivedCode(ctx context.Context, order core.Order, previousStatus core.OrderStatus, result core.ProviderCodeResult) (core.Order, *core.SMSCode, error) {
	receivedAt := result.ReceivedAt
	if receivedAt.IsZero() {
		receivedAt = order.UpdatedAt
	}
	code := core.SMSCode{Value: result.Code, MessageText: result.MessageText, ReceivedAt: receivedAt}
	code, err := s.prepareCodeSecret(ctx, order, code)
	if err != nil {
		return core.Order{}, nil, err
	}
	order.Status = core.StatusCodeReceived
	records, err := s.statusAndCodeRecords(ctx, order, previousStatus, code)
	if err != nil {
		return core.Order{}, nil, err
	}
	if err := s.recordCode(ctx, order, code, records...); err != nil {
		return core.Order{}, nil, err
	}
	return order, &code, nil
}
