package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) newAcquireRequestOrder(ctx context.Context, cmd core.AcquireNumberCommand, route core.Route) core.Order {
	requestID := firstNonEmpty(cmd.RequestID, s.ids.NewID("req_"))
	now := s.clock.Now()
	return core.Order{
		ID:          s.ids.NewID("ord_"),
		RequestID:   requestID,
		ProviderKey: route.ProviderKey,
		Target:      withRouteTargetDefaults(core.Target{}, route),
		Status:      core.StatusAcquireRequested,
		ExpiresAt:   orderRequestExpiresAt(now, s.routePolicy(ctx, route), cmd.LeaseDuration),
		UpdatedAt:   now,
	}
}
