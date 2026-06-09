package eventbusadapter

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/eventoutbox"
)

func (b *OrderEventRecorder) CodeReceived(ctx context.Context, order core.Order, code core.SMSCode) (eventoutbox.Record, error) {
	return b.recordCodeReceived(ctx, order, code)
}

func (b *OrderEventRecorder) OrderStatusChanged(ctx context.Context, order core.Order, previous core.OrderStatus) (eventoutbox.Record, error) {
	return b.recordOrderStatusChanged(ctx, order, previous)
}
