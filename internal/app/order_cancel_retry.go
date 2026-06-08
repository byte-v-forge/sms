package app

import (
	"context"
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func shouldQueueEarlyCancelRetry(err *core.Error, policy core.ProviderPolicy) bool {
	return err != nil && err.Code == core.CodeCancelNotAllowed && err.Retryable && policy.EarlyCancelRetryAfter > 0
}

func earlyCancelRetryAt(order core.Order, policy core.ProviderPolicy, now time.Time) time.Time {
	if !order.AcquiredAt.IsZero() && policy.EarlyCancelRetryAfter > 0 {
		retryAt := order.AcquiredAt.Add(policy.EarlyCancelRetryAfter)
		if retryAt.After(now) {
			return retryAt
		}
	}
	delay := policy.PollInterval
	if delay <= 0 {
		delay = 5 * time.Second
	}
	return now.Add(delay)
}

func (s *OrderService) queueCancelRequest(ctx context.Context, order core.Order, requestID string, reason string) (core.Order, error) {
	if requestID = strings.TrimSpace(requestID); requestID == "" {
		requestID = s.ids.NewID("req_")
	}
	retryAt := s.cancelReadyAt(ctx, order)
	order = markCancelRequest(order, requestID, retryAt, s.clock.Now())
	record, err := s.events.OrderCancelRequested(ctx, order, requestID, reason)
	if err != nil {
		return order, err
	}
	if err := s.updateOrder(ctx, order, record); err != nil {
		return core.Order{}, err
	}
	return s.execution.AfterCancelQueued(ctx, s, order, requestID)
}

func (s *OrderService) deferCancelRetry(ctx context.Context, order core.Order, requestID string, retryAt time.Time) (core.Order, error) {
	if retryAt.IsZero() {
		retryAt = s.clock.Now().Add(5 * time.Second)
	}
	if !retryAt.After(s.clock.Now()) {
		retryAt = s.clock.Now().Add(5 * time.Second)
	}
	order = markCancelRequest(order, requestID, retryAt, s.clock.Now())
	if err := s.updateOrder(ctx, order); err != nil {
		return core.Order{}, err
	}
	return order, &CancelRetryError{RetryAt: retryAt}
}

func (s *OrderService) cancelReadyAt(ctx context.Context, order core.Order) time.Time {
	now := s.clock.Now()
	provider, err := s.provider(order.ProviderKey)
	if err != nil {
		return now
	}
	policy := providerPolicyForOrder(ctx, provider, order).WithDefaults()
	if policy.CancelAllowedAfter <= 0 || order.AcquiredAt.IsZero() {
		return now
	}
	readyAt := order.AcquiredAt.Add(policy.CancelAllowedAfter)
	if readyAt.After(now) {
		return readyAt
	}
	return now
}
