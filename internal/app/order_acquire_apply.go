package app

import (
	"context"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) applyProviderOrder(ctx context.Context, order core.Order, route core.Route, provider core.Provider, providerOrder core.ProviderOrder) (core.Order, error) {
	now := s.clock.Now()
	acquiredAt := providerOrder.AcquiredAt
	if acquiredAt.IsZero() {
		acquiredAt = now
	}
	policy := providerPolicyForUpstreamOrder(provider, providerOrder.UpstreamOrderID).WithDefaults()
	expiresAt := providerOrder.ExpiresAt
	if expiresAt.IsZero() {
		expiresAt = acquiredAt.Add(policy.OrderTTL)
	}
	if !order.ExpiresAt.IsZero() && order.ExpiresAt.Before(expiresAt) {
		expiresAt = order.ExpiresAt
	}
	var cancelAllowedAt time.Time
	if policy.CancelAllowedAfter > 0 {
		cancelAllowedAt = acquiredAt.Add(policy.CancelAllowedAfter)
	}
	previousStatus := order.Status
	order.ProviderKey = provider.Key()
	order.UpstreamOrderID = providerOrder.UpstreamOrderID
	order.Target = withRouteTargetDefaults(order.Target, route)
	order.PhoneNumber = providerOrder.PhoneNumber
	order.Status = core.StatusPendingCode
	order.Price = providerOrder.Price
	order.AcquiredAt = acquiredAt
	order.ExpiresAt = expiresAt
	order.UpdatedAt = now
	order.CancelAllowedAt = cancelAllowedAt
	order.CanRequestAdditionalCode = providerOrder.CanRequestAdditionalCode
	order.LastError = nil
	records, err := s.orderAcquiredRecords(ctx, order, "order_acquired")
	if err != nil {
		return core.Order{}, err
	}
	statusRecords, err := s.statusChangedRecords(ctx, order, previousStatus)
	if err != nil {
		return core.Order{}, err
	}
	records = append(records, statusRecords...)
	if err := s.updateOrder(ctx, order, records...); err != nil {
		return core.Order{}, err
	}
	return order, nil
}

func (s *OrderService) recordAcquireFailure(ctx context.Context, order core.Order, smsErr *core.Error) (core.Order, error) {
	now := s.clock.Now()
	order.LastError = smsErr
	order.UpdatedAt = now
	if smsErr.Retryable {
		if err := s.updateOrder(ctx, order); err != nil {
			return core.Order{}, err
		}
		return order, smsErr
	}
	previousStatus := order.Status
	order.Status = core.StatusFailed
	records, err := s.statusChangedRecords(ctx, order, previousStatus)
	if err != nil {
		return core.Order{}, err
	}
	if err := s.updateOrder(ctx, order, records...); err != nil {
		return core.Order{}, err
	}
	return order, smsErr
}
