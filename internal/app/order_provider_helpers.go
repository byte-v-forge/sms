package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) provider(key string) (core.Provider, error) {
	provider, ok := s.providers[key]
	if !ok {
		return nil, core.NewError(core.CodeRouteNotFound, "sms provider not registered", false)
	}
	return provider, nil
}

func (s *OrderService) routePolicy(ctx context.Context, route core.Route) core.ProviderPolicy {
	provider, err := s.provider(route.ProviderKey)
	if err != nil {
		return core.ProviderPolicy{}.WithDefaults()
	}
	if configured, ok := provider.(*ConfiguredProvider); ok {
		return configured.LoadPolicyForOrder(ctx, "").WithDefaults()
	}
	return provider.Policy().WithDefaults()
}
