package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) CheckCode(ctx context.Context, orderID string) (core.Order, *core.SMSCode, error) {
	order, err := s.store.Get(ctx, orderID)
	if err != nil {
		return core.Order{}, nil, err
	}
	if err := s.expireIfNeeded(ctx, &order); err != nil {
		return core.Order{}, nil, err
	}
	if handled, err := orderCodeCheckPrecondition(order); handled {
		return order, nil, err
	}
	previousStatus := order.Status
	provider, err := s.provider(order.ProviderKey)
	if err != nil {
		return core.Order{}, nil, err
	}
	bindOrderProviderConfig(provider, order)
	result, err := s.providerStatus(ctx, order, provider)
	if err != nil {
		return order, nil, err
	}
	order.UpdatedAt = s.clock.Now()
	if result.Status == core.StatusCodeReceived {
		return s.applyReceivedCode(ctx, order, previousStatus, result)
	}
	return s.applyProviderStatus(ctx, order, previousStatus, result)
}
