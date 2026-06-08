package app

import (
	"context"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) handleProviderCancelError(ctx context.Context, order core.Order, provider core.Provider, policy core.ProviderPolicy, requestID string, now time.Time, err error) (core.Order, error) {
	smsErr := asCoreError(err)
	if raced, ok, raceErr := s.resolveCancelRace(ctx, order, provider, smsErr); raceErr != nil {
		return order, raceErr
	} else if ok {
		return raced, nil
	}
	if shouldQueueEarlyCancelRetry(smsErr, policy) {
		return s.deferCancelRetry(ctx, order, requestID, earlyCancelRetryAt(order, policy, now))
	}
	order.LastError = smsErr
	order.UpdatedAt = now
	if smsErr.Retryable {
		order.CancelAllowedAt = earlyCancelRetryAt(order, policy, now)
	}
	_ = s.updateOrder(ctx, order)
	return order, err
}
