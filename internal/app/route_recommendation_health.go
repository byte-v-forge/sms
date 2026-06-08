package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *CatalogService) disabledRouteKeys(ctx context.Context, offers []core.RouteOffer, failurePolicy core.RouteFailurePolicy) map[string]struct{} {
	if len(offers) == 0 || s.routeHealth == nil {
		return map[string]struct{}{}
	}
	routes := make([]core.Route, 0, len(offers))
	for _, offer := range offers {
		routes = append(routes, routeWithFailurePolicy(offer.Route, failurePolicy))
	}
	disabled, err := s.routeHealth.DisabledRouteKeys(ctx, routes)
	if err != nil {
		return map[string]struct{}{}
	}
	return disabled
}

func routeWithFailurePolicy(route core.Route, failurePolicy core.RouteFailurePolicy) core.Route {
	route.FailurePolicy = failurePolicy
	return route
}

func routeTemporarilyDisabled(route core.Route, disabledRoutes map[string]struct{}) bool {
	if len(disabledRoutes) == 0 {
		return false
	}
	_, disabled := disabledRoutes[routeHealthKey(route)]
	return disabled
}
