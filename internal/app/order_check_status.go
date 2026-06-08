package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) applyProviderStatus(ctx context.Context, order core.Order, previousStatus core.OrderStatus, result core.ProviderCodeResult) (core.Order, *core.SMSCode, error) {
	switch result.Status {
	case core.StatusPendingCode, core.StatusAdditionalCodeRequested:
		order.Status = result.Status
	case core.StatusCanceled, core.StatusFailed, core.StatusExpired:
		order.Status = result.Status
	}
	records, err := s.statusChangedRecords(ctx, order, previousStatus)
	if err != nil {
		return core.Order{}, nil, err
	}
	if err := s.updateOrder(ctx, order, records...); err != nil {
		return core.Order{}, nil, err
	}
	return order, nil, nil
}
