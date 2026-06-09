package eventbusadapter

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/eventoutbox"
)

func (b *OrderEventRecorder) OrderAcquireRequested(ctx context.Context, order core.Order, route core.Route, reason string) (eventoutbox.Record, error) {
	return b.recordOrderAcquireRequested(ctx, order, route, reason)
}

func (b *OrderEventRecorder) OrderAcquired(ctx context.Context, order core.Order) (eventoutbox.Record, error) {
	return b.recordOrderAcquired(ctx, order)
}
