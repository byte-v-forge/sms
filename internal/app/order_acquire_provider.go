package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) acquireWithRoute(ctx context.Context, order core.Order, requestID string, route core.Route) (core.Order, error) {
	provider, err := s.provider(route.ProviderKey)
	if err != nil {
		return core.Order{}, err
	}
	providerOrder, err := provider.AcquireNumber(ctx, s.providerAcquireRequest(order, requestID, route))
	if err != nil {
		return core.Order{}, err
	}
	_ = s.routeHealth.RecordAcquireSuccess(ctx, route)
	return s.applyProviderOrder(ctx, order, route, provider, providerOrder)
}

func (s *OrderService) providerAcquireRequest(order core.Order, requestID string, route core.Route) core.ProviderAcquireRequest {
	return core.ProviderAcquireRequest{
		RequestID:     firstNonEmpty(requestID, order.RequestID),
		Route:         route,
		Target:        withRouteTargetDefaults(order.Target, route),
		LeaseDuration: remainingLease(s.clock.Now(), order.ExpiresAt),
	}
}
