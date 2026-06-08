package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (p *ConfiguredProvider) AcquireNumber(ctx context.Context, request core.ProviderAcquireRequest) (core.ProviderOrder, error) {
	provider, err := p.providerForTarget(ctx, request.Target)
	if err != nil {
		return core.ProviderOrder{}, err
	}
	order, err := provider.AcquireNumber(ctx, request)
	if err == nil && order.UpstreamOrderID != "" {
		order = withProviderOrderExpiry(order, provider.Policy().WithDefaults())
	}
	return order, err
}

func withProviderOrderExpiry(order core.ProviderOrder, policy core.ProviderPolicy) core.ProviderOrder {
	if order.ExpiresAt.IsZero() && !order.AcquiredAt.IsZero() && policy.OrderTTL > 0 {
		order.ExpiresAt = order.AcquiredAt.Add(policy.OrderTTL)
	}
	return order
}
