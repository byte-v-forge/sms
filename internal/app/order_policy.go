package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func bindOrderProviderConfig(provider core.Provider, order core.Order) {
	configured, ok := provider.(*ConfiguredProvider)
	if !ok {
		return
	}
	configured.BindOrderConfig(order.UpstreamOrderID)
}

func providerPolicyForOrder(ctx context.Context, provider core.Provider, order core.Order) core.ProviderPolicy {
	if configured, ok := provider.(*ConfiguredProvider); ok {
		return configured.LoadPolicyForOrder(ctx, order.UpstreamOrderID)
	}
	return providerPolicyForUpstreamOrder(provider, order.UpstreamOrderID)
}

func providerPolicyForUpstreamOrder(provider core.Provider, upstreamOrderID string) core.ProviderPolicy {
	if configured, ok := provider.(*ConfiguredProvider); ok {
		return configured.PolicyForOrder(upstreamOrderID)
	}
	return provider.Policy()
}
