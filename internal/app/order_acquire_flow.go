package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) RunAcquireRequest(ctx context.Context, orderID string, requestID string, route core.Route) (core.Order, error) {
	order, route, ok, err := s.prepareAcquireRequest(ctx, orderID, route)
	if err != nil || !ok {
		return order, err
	}
	acquired, acquireErr := s.acquireWithRoute(ctx, order, requestID, route)
	if acquireErr == nil {
		return acquired, nil
	}
	smsErr := asCoreError(acquireErr)
	_ = s.routeHealth.RecordAcquireFailure(ctx, route, smsErr)
	return s.recordAcquireFailure(ctx, order, smsErr)
}
