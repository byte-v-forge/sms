package app

import (
	"context"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) PollInterval(ctx context.Context, order core.Order) time.Duration {
	provider, err := s.provider(order.ProviderKey)
	if err != nil {
		return 5 * time.Second
	}
	return providerPolicyForOrder(ctx, provider, order).WithDefaults().PollInterval
}
