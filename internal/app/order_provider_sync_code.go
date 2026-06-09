package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) applySyncedProviderCode(
	ctx context.Context,
	order core.Order,
	previousStatus core.OrderStatus,
	result core.ProviderCodeResult,
) (core.Order, error) {
	code := codeFromProviderResult(result, order.UpdatedAt)
	code, err := s.prepareCodeSecret(ctx, order, code)
	if err != nil {
		return order, err
	}
	order.Status = core.StatusCodeReceived
	records, err := s.statusAndCodeRecords(ctx, order, previousStatus, code)
	if err != nil {
		return order, err
	}
	if err := s.recordCode(ctx, order, code, records...); err != nil {
		return core.Order{}, err
	}
	return order, nil
}
