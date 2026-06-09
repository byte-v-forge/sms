package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) applySyncedProviderResult(ctx context.Context, order core.Order, previousStatus core.OrderStatus, result core.ProviderCodeResult) (core.Order, bool, error) {
	switch result.Status {
	case core.StatusCodeReceived:
		synced, err := s.applySyncedProviderCode(ctx, order, previousStatus, result)
		return synced, true, err
	case core.StatusCompleted, core.StatusCanceled, core.StatusExpired, core.StatusFailed:
		synced, err := s.applySyncedProviderTerminalStatus(ctx, order, previousStatus, result.Status)
		return synced, true, err
	default:
		return order, false, nil
	}
}
