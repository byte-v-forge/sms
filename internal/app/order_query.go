package app

import (
	"context"
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *OrderService) GetOrder(ctx context.Context, orderID string) (core.Order, error) {
	order, err := s.store.Get(ctx, orderID)
	if err != nil {
		return core.Order{}, err
	}
	if err := s.expireIfNeeded(ctx, &order); err != nil {
		return core.Order{}, err
	}
	if s.execution.SyncProviderStateOnRead() && order.Status.HasProviderLease() && order.Status != core.StatusCodeReceived && !order.Status.IsFinal() {
		synced, _, err := s.CheckCode(ctx, orderID)
		if synced.ID == "" {
			synced = order
		}
		if err == nil {
			return synced, nil
		}
		return synced, err
	}
	return order, nil
}

func (s *OrderService) CheckCode(ctx context.Context, orderID string) (core.Order, *core.SMSCode, error) {
	order, err := s.store.Get(ctx, orderID)
	if err != nil {
		return core.Order{}, nil, err
	}
	if err := s.expireIfNeeded(ctx, &order); err != nil {
		return core.Order{}, nil, err
	}
	if order.Status == core.StatusCodeReceived {
		return order, nil, nil
	}
	if order.Status == core.StatusAcquireRequested || strings.TrimSpace(order.UpstreamOrderID) == "" {
		return order, nil, nil
	}
	if order.Status.IsFinal() {
		return order, nil, core.NewError(core.CodeOrderAlreadyFinalized, "order already finalized", false)
	}
	previousStatus := order.Status
	provider, err := s.provider(order.ProviderKey)
	if err != nil {
		return core.Order{}, nil, err
	}
	bindOrderProviderConfig(provider, order)
	result, err := provider.GetStatus(ctx, order.UpstreamOrderID)
	if err != nil {
		order.LastError = asCoreError(err)
		order.UpdatedAt = s.clock.Now()
		_ = s.updateOrder(ctx, order)
		return order, nil, err
	}
	order.UpdatedAt = s.clock.Now()
	switch result.Status {
	case core.StatusCodeReceived:
		receivedAt := result.ReceivedAt
		if receivedAt.IsZero() {
			receivedAt = order.UpdatedAt
		}
		code := core.SMSCode{Value: result.Code, MessageText: result.MessageText, ReceivedAt: receivedAt}
		code, err = s.prepareCodeSecret(ctx, order, code)
		if err != nil {
			return core.Order{}, nil, err
		}
		order.Status = core.StatusCodeReceived
		records, err := s.statusAndCodeRecords(ctx, order, previousStatus, code)
		if err != nil {
			return core.Order{}, nil, err
		}
		if err := s.recordCode(ctx, order, code, records...); err != nil {
			return core.Order{}, nil, err
		}
		return order, &code, nil
	case core.StatusPendingCode, core.StatusAdditionalCodeRequested:
		order.Status = result.Status
	case core.StatusCanceled, core.StatusFailed, core.StatusExpired:
		order.Status = result.Status
	}
	records, err := s.statusChangedRecords(ctx, order, previousStatus)
	if err != nil {
		return core.Order{}, nil, err
	}
	if err := s.updateOrder(ctx, order, records...); err != nil {
		return core.Order{}, nil, err
	}
	return order, nil, nil
}

func (s *OrderService) PollInterval(ctx context.Context, order core.Order) time.Duration {
	provider, err := s.provider(order.ProviderKey)
	if err != nil {
		return 5 * time.Second
	}
	return providerPolicyForOrder(ctx, provider, order).WithDefaults().PollInterval
}
