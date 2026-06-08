package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

type OrderExecution interface {
	AfterAcquireRequested(context.Context, *OrderService, core.Order, core.Route) (core.Order, error)
	AfterCancelQueued(context.Context, *OrderService, core.Order, string) (core.Order, error)
	SyncProviderStateOnRead() bool
}

type inProcessOrderExecution struct{}

func (inProcessOrderExecution) AfterAcquireRequested(ctx context.Context, service *OrderService, order core.Order, route core.Route) (core.Order, error) {
	return service.RunAcquireRequest(ctx, order.ID, order.RequestID, route)
}

func (inProcessOrderExecution) AfterCancelQueued(ctx context.Context, service *OrderService, order core.Order, requestID string) (core.Order, error) {
	return service.cancelLoadedOrder(ctx, order, requestID)
}

func (inProcessOrderExecution) SyncProviderStateOnRead() bool {
	return true
}

type asyncOrderExecution struct{}

func (asyncOrderExecution) AfterAcquireRequested(_ context.Context, _ *OrderService, order core.Order, _ core.Route) (core.Order, error) {
	return order, nil
}

func (asyncOrderExecution) AfterCancelQueued(_ context.Context, _ *OrderService, order core.Order, _ string) (core.Order, error) {
	return order, nil
}

func (asyncOrderExecution) SyncProviderStateOnRead() bool {
	return false
}

type asyncOrderEventSink interface {
	AsyncRequests() bool
}

func orderExecutionForEvents(events OrderEventSink) OrderExecution {
	async, ok := events.(asyncOrderEventSink)
	if ok && async.AsyncRequests() {
		return asyncOrderExecution{}
	}
	return inProcessOrderExecution{}
}
