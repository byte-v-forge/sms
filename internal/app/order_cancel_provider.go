package app

import (
	"context"
	"strings"

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
		return s.expireLoadedOrder(ctx, order, now)
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
	return s.completeProviderCancel(ctx, order, now)
}
