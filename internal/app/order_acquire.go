package app

import (
	"context"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) AcquireNumber(ctx context.Context, cmd core.AcquireNumberCommand) (core.Order, error) {
	route := cmd.AcquireParams
	if err := validateAcquireRoute(route); err != nil {
		return core.Order{}, err
	}
	if cmd.RequestID == "" {
		cmd.RequestID = s.ids.NewID("req_")
	}
	now := s.clock.Now()
	target := withRouteTargetDefaults(core.Target{}, route)
	expiresAt := orderRequestExpiresAt(now, s.routePolicy(ctx, route), cmd.LeaseDuration)
	order := core.Order{
		ID:          s.ids.NewID("ord_"),
		RequestID:   cmd.RequestID,
		ProviderKey: route.ProviderKey,
		Target:      target,
		Status:      core.StatusAcquireRequested,
		ExpiresAt:   expiresAt,
		UpdatedAt:   now,
	}
	record, err := s.events.OrderAcquireRequested(ctx, order, route, "api_request")
	if err != nil {
		return core.Order{}, err
	}
	if err := s.saveOrder(ctx, order, record); err != nil {
		return core.Order{}, err
	}
	return s.execution.AfterAcquireRequested(ctx, s, order, route)
}

func validateAcquireRoute(route core.Route) error {
	if strings.TrimSpace(route.ProviderKey) == "" || strings.TrimSpace(route.UpstreamServiceKey) == "" || strings.TrimSpace(route.ProviderCountryID) == "" {
		return core.NewError(core.CodeValidationFailed, "sms acquire params are incomplete", false)
	}
	switch normalizeProviderKey(route.ProviderKey) {
	case "5sim", "smsbower":
		if strings.TrimSpace(route.UpstreamProviderID) == "" {
			return core.NewError(core.CodeValidationFailed, "sms upstream provider id is required", false)
		}
	}
	return nil
}
