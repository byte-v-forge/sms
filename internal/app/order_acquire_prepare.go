package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) prepareAcquireRequest(ctx context.Context, orderID string, route core.Route) (core.Order, core.Route, bool, error) {
	order, err := s.store.Get(ctx, orderID)
	if err != nil {
		return core.Order{}, core.Route{}, false, err
	}
	if order.Status != core.StatusAcquireRequested {
		return order, route, false, nil
	}
	if err := s.expireIfNeeded(ctx, &order); err != nil {
		return order, route, false, err
	}
	route = acquireRequestRoute(order, route)
	if route.ProviderKey == "" {
		failure := core.NewError(core.CodeRouteNotFound, "sms acquire route not found", false)
		failed, err := s.recordAcquireFailure(ctx, order, failure)
		return failed, route, false, err
	}
	if err := validateAcquireRoute(route); err != nil {
		failed, err := s.recordAcquireFailure(ctx, order, asCoreError(err))
		return failed, route, false, err
	}
	return order, route, true, nil
}
