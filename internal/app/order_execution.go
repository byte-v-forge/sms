package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

type OrderExecution interface {
	AfterAcquireRequested(context.Context, core.Order, core.Route) (core.Order, error)
	AfterCancelQueued(context.Context, core.Order, string) (core.Order, error)
	SyncProviderStateOnRead() bool
}

type acquireRequestRunner func(context.Context, string, string, core.Route) (core.Order, error)

type cancelRequestRunner func(context.Context, core.Order, string) (core.Order, error)

type inProcessOrderExecution struct {
	acquire acquireRequestRunner
	cancel  cancelRequestRunner
}

func (e inProcessOrderExecution) AfterAcquireRequested(ctx context.Context, order core.Order, route core.Route) (core.Order, error) {
	return e.acquire(ctx, order.ID, order.RequestID, route)
}

func (e inProcessOrderExecution) AfterCancelQueued(ctx context.Context, order core.Order, requestID string) (core.Order, error) {
	return e.cancel(ctx, order, requestID)
}

func (inProcessOrderExecution) SyncProviderStateOnRead() bool {
	return true
}

type asyncOrderExecution struct{}

func (asyncOrderExecution) AfterAcquireRequested(_ context.Context, order core.Order, _ core.Route) (core.Order, error) {
	return order, nil
}

func (asyncOrderExecution) AfterCancelQueued(_ context.Context, order core.Order, _ string) (core.Order, error) {
	return order, nil
}

func (asyncOrderExecution) SyncProviderStateOnRead() bool {
	return false
}

type asyncOrderEventSink interface {
	AsyncRequests() bool
}

func orderExecutionForEvents(events OrderEventSink, acquire acquireRequestRunner, cancel cancelRequestRunner) OrderExecution {
	async, ok := events.(asyncOrderEventSink)
	if ok && async.AsyncRequests() {
		return asyncOrderExecution{}
	}
	return inProcessOrderExecution{acquire: acquire, cancel: cancel}
}
