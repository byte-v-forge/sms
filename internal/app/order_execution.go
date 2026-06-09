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
