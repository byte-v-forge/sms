package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

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
