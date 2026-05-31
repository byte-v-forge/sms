package app

import (
	"context"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) RunAcquireRequest(ctx context.Context, orderID string, requestID string, route core.Route) (core.Order, error) {
	order, err := s.store.Get(ctx, orderID)
	if err != nil {
		return core.Order{}, err
	}
	if order.Status != core.StatusAcquireRequested {
		return order, nil
	}
	if err := s.expireIfNeeded(ctx, &order); err != nil {
		return order, err
	}
	if strings.TrimSpace(route.ProviderKey) == "" {
		route = routeFromOrder(order)
	}
	if strings.TrimSpace(route.ProviderKey) == "" {
		return s.recordAcquireFailure(ctx, order, core.NewError(core.CodeRouteNotFound, "sms acquire route not found", false))
	}
	if err := validateAcquireRoute(route); err != nil {
		return s.recordAcquireFailure(ctx, order, asCoreError(err))
	}
	acquired, acquireErr := s.acquireWithRoute(ctx, order, requestID, route)
	if acquireErr == nil {
		return acquired, nil
	}
	smsErr := asCoreError(acquireErr)
	_ = s.routeHealth.RecordAcquireFailure(ctx, route, smsErr)
	return s.recordAcquireFailure(ctx, order, smsErr)
}

func (s *OrderService) acquireWithRoute(ctx context.Context, order core.Order, requestID string, route core.Route) (core.Order, error) {
	provider, err := s.provider(route.ProviderKey)
	if err != nil {
		return core.Order{}, err
	}
	target := withRouteTargetDefaults(order.Target, route)
	providerOrder, err := provider.AcquireNumber(ctx, core.ProviderAcquireRequest{
		RequestID:     firstNonEmpty(requestID, order.RequestID),
		Route:         route,
		Target:        target,
		LeaseDuration: remainingLease(s.clock.Now(), order.ExpiresAt),
	})
	if err != nil {
		return core.Order{}, err
	}
	_ = s.routeHealth.RecordAcquireSuccess(ctx, route)
	return s.applyProviderOrder(ctx, order, route, provider, providerOrder)
}
