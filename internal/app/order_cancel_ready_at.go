package app

import (
	"context"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

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
