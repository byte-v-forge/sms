package app

import (
	"context"
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) cancelLoadedOrder(ctx context.Context, order core.Order, requestID string) (core.Order, error) {
	provider, err := s.provider(order.ProviderKey)
	if err != nil {
		return core.Order{}, err
	}
	bindOrderProviderConfig(provider, order)
	policy := providerPolicyForOrder(ctx, provider, order).WithDefaults()
	now := s.clock.Now()
	if order.Status.IsFinal() {
		return order, core.NewError(core.CodeOrderAlreadyFinalized, "order already finalized", false)
	}
	if orderHasCode(order) {
		return order, nil
	}
	if !order.Status.HasProviderLease() || strings.TrimSpace(order.UpstreamOrderID) == "" {
		return s.cancelLocalOrder(ctx, order, now, true)
	}
	if order.IsExpired(now) {
		previousStatus := order.Status
		order.Status = core.StatusExpired
		order.UpdatedAt = now
		records, err := s.statusChangedRecords(ctx, order, previousStatus)
		if err != nil {
			return order, err
		}
		if err := s.updateOrder(ctx, order, records...); err != nil {
			return order, err
		}
		return order, core.NewError(core.CodeOrderExpired, "order expired", false)
	}
	if synced, handled, syncErr := s.syncTerminalOrChargedProviderState(ctx, order, provider); syncErr != nil {
		return order, syncErr
	} else if handled {
		return synced, nil
	}
	order = normalizeProviderCancelTimes(order, policy, now)
	if !order.CancelAllowedAt.IsZero() && now.Before(order.CancelAllowedAt) {
		return order, &CancelRetryError{RetryAt: order.CancelAllowedAt}
	}
	age := now.Sub(order.AcquiredAt)
	if policy.CancelAllowedAfter > 0 && age < policy.CancelAllowedAfter {
		return s.deferCancelRetry(ctx, order, requestID, order.AcquiredAt.Add(policy.CancelAllowedAfter))
	}
	if policy.CancelAllowedUntil > 0 && age > policy.CancelAllowedUntil {
		return order, core.NewError(core.CodeCancelNotAllowed, "order is too old to cancel", false)
	}
	if err := provider.SetStatus(ctx, order.UpstreamOrderID, core.ActionCancelOrder); err != nil {
		return s.handleProviderCancelError(ctx, order, provider, policy, requestID, now, err)
	}
	previousStatus := order.Status
	order.Status = core.StatusCanceled
	order.UpdatedAt = now
	order.LastError = nil
	order.CancelAllowedAt = time.Time{}
	records, err := s.statusChangedRecords(ctx, order, previousStatus)
	if err != nil {
		return core.Order{}, err
	}
	if err := s.updateOrder(ctx, order, records...); err != nil {
		return core.Order{}, err
	}
	return order, nil
}

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
