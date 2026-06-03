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
	now := s.clock.Now()
	order.UpdatedAt = now
	order.LastError = nil
	order.CancelAllowedAt = time.Time{}
	switch result.Status {
	case core.StatusCodeReceived:
		receivedAt := result.ReceivedAt
		if receivedAt.IsZero() {
			receivedAt = now
		}
		code := core.SMSCode{Value: result.Code, MessageText: result.MessageText, ReceivedAt: receivedAt}
		code, err := s.prepareCodeSecret(ctx, order, code)
		if err != nil {
			return order, true, err
		}
		order.Status = core.StatusCodeReceived
		records, err := s.statusAndCodeRecords(ctx, order, previousStatus, code)
		if err != nil {
			return order, true, err
		}
		if err := s.recordCode(ctx, order, code, records...); err != nil {
			return core.Order{}, true, err
		}
		return order, true, nil
	case core.StatusCompleted, core.StatusCanceled, core.StatusExpired, core.StatusFailed:
		order.Status = result.Status
		records, err := s.statusChangedRecords(ctx, order, previousStatus)
		if err != nil {
			return order, true, err
		}
		if err := s.updateOrder(ctx, order, records...); err != nil {
			return core.Order{}, true, err
		}
		return order, true, nil
	default:
		return order, false, nil
	}
}
