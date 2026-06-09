package app

import (
	"context"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) resolveCancelRace(ctx context.Context, order core.Order, provider core.Provider, cancelErr *core.Error) (core.Order, bool, error) {
	if cancelErr == nil {
		return order, false, nil
	}
	if cancelErr.Retryable && cancelErr.Code != core.CodeCancelNotAllowed {
		return order, false, nil
	}
	synced, handled, err := s.syncTerminalOrChargedProviderState(ctx, order, provider)
	if err != nil {
		return order, false, nil
	}
	return synced, handled, nil
}

func (s *OrderService) syncTerminalOrChargedProviderState(ctx context.Context, order core.Order, provider core.Provider) (core.Order, bool, error) {
	result, err := provider.GetStatus(ctx, order.UpstreamOrderID)
	if err != nil {
		return order, false, err
	}
	previousStatus := order.Status
	order = prepareSyncedProviderOrder(order, s.clock.Now())
	switch result.Status {
	case core.StatusCodeReceived:
		order, err = s.applySyncedProviderCode(ctx, order, previousStatus, result)
		return order, true, err
	case core.StatusCompleted, core.StatusCanceled, core.StatusExpired, core.StatusFailed:
		order, err = s.applySyncedProviderTerminalStatus(ctx, order, previousStatus, result.Status)
		return order, true, err
	default:
		return order, false, nil
	}
}

func prepareSyncedProviderOrder(order core.Order, now time.Time) core.Order {
	order.UpdatedAt = now
	order.LastError = nil
	order.CancelAllowedAt = time.Time{}
	return order
}
