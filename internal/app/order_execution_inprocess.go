package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

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
