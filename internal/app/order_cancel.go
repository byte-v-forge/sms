package app

import (
	"context"
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

type CancelRetryError struct {
	RetryAt time.Time
}

func (e *CancelRetryError) Error() string {
	if e == nil || e.RetryAt.IsZero() {
		return "sms order cancel retry is scheduled"
	}
	return "sms order cancel retry is scheduled for " + e.RetryAt.UTC().Format(time.RFC3339)
}

func (s *OrderService) CancelOrder(ctx context.Context, orderID, requestID string) (core.Order, error) {
	order, err := s.store.Get(ctx, orderID)
	if err != nil {
		return core.Order{}, err
	}
	if err := s.expireIfNeeded(ctx, &order); err != nil {
		return order, err
	}
	if order.Status.IsFinal() {
		return order, core.NewError(core.CodeOrderAlreadyFinalized, "order already finalized", false)
	}
	if orderHasCode(order) {
		return order, nil
	}
	if !order.Status.HasProviderLease() || strings.TrimSpace(order.UpstreamOrderID) == "" {
		previousStatus := order.Status
		order.Status = core.StatusCanceled
		order.UpdatedAt = s.clock.Now()
		order.LastError = nil
		records, err := s.statusChangedRecords(ctx, order, previousStatus)
		if err != nil {
			return order, err
		}
		if err := s.updateOrder(ctx, order, records...); err != nil {
			return core.Order{}, err
		}
		return order, nil
	}
	return s.queueCancelRequest(ctx, order, requestID, "api_request")
}

func (s *OrderService) RunCancelRequest(ctx context.Context, orderID, requestID string) (core.Order, error) {
	order, err := s.store.Get(ctx, orderID)
	if err != nil {
		return core.Order{}, err
	}
	return s.cancelLoadedOrder(ctx, order, requestID)
}

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
		previousStatus := order.Status
		order.Status = core.StatusCanceled
		order.UpdatedAt = now
		order.LastError = nil
		order.CancelAllowedAt = time.Time{}
		records, err := s.statusChangedRecords(ctx, order, previousStatus)
		if err != nil {
			return order, err
		}
		if err := s.updateOrder(ctx, order, records...); err != nil {
			return core.Order{}, err
		}
		return order, nil
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

func normalizeProviderCancelTimes(order core.Order, policy core.ProviderPolicy, now time.Time) core.Order {
	if order.AcquiredAt.After(now.Add(providerClockSkewTolerance)) {
		order.AcquiredAt = now
		if policy.CancelAllowedAfter > 0 {
			order.AcquiredAt = now.Add(-policy.CancelAllowedAfter)
		}
	}
	if !order.CancelAllowedAt.IsZero() && !order.ExpiresAt.IsZero() && order.CancelAllowedAt.After(order.ExpiresAt) {
		order.CancelAllowedAt = time.Time{}
	}
	if !order.CancelAllowedAt.IsZero() && order.CancelAllowedAt.After(now.Add(providerClockSkewTolerance)) {
		order.CancelAllowedAt = order.AcquiredAt.Add(policy.CancelAllowedAfter)
	}
	return order
}
