package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) AcquireNumber(ctx context.Context, cmd core.AcquireNumberCommand) (core.Order, error) {
	route := cmd.AcquireParams
	if err := validateAcquireRoute(route); err != nil {
		return core.Order{}, err
	}
	order := s.newAcquireRequestOrder(ctx, cmd, route)
	record, err := s.events.OrderAcquireRequested(ctx, order, route, "api_request")
	if err != nil {
		return core.Order{}, err
	}
	if err := s.saveOrder(ctx, order, record); err != nil {
		return core.Order{}, err
	}
	return s.execution.AfterAcquireRequested(ctx, order, route)
}
