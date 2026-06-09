package eventbusadapter

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/eventoutbox"
)

func (b *OrderEventRecorder) OrderPollRequested(ctx context.Context, order core.Order, reason string) (eventoutbox.Record, error) {
	return b.recordOrderPollRequested(ctx, order, reason)
}

func (b *OrderEventRecorder) OrderCancelRequested(ctx context.Context, order core.Order, requestID string, reason string) (eventoutbox.Record, error) {
	return b.recordOrderCancelRequested(ctx, order, requestID, reason)
}
