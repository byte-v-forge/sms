package app

import (
	"context"

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
	if err := validateCancelableOrder(order); err != nil {
		return order, err
	}
	if orderHasCode(order) {
		return order, nil
	}
	if orderHasNoProviderLease(order) {
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
	return s.cancelActiveProviderOrder(ctx, order, provider, policy, requestID, now)
}
