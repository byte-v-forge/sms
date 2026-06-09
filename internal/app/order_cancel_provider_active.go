package app

import (
	"context"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) cancelActiveProviderOrder(ctx context.Context, order core.Order, provider core.Provider, policy core.ProviderPolicy, requestID string, now time.Time) (core.Order, error) {
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
	return s.completeProviderCancel(ctx, order, now)
}
