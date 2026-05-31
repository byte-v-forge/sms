package eventbusadapter

import (
	"context"
	"fmt"

	"github.com/byte-v-forge/common-lib/eventbus"
	"github.com/byte-v-forge/common-lib/eventcatalog"
	"github.com/byte-v-forge/common-lib/eventoutbox"
	smsv1 "github.com/byte-v-forge/common-lib/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/app"
	"github.com/byte-v-forge/sms/internal/core"
)

func (b *OrderEventRecorder) CodeReceived(ctx context.Context, order core.Order, code core.SMSCode) (eventoutbox.Record, error) {
	eventCtx := b.context(
		eventcatalog.SMSCodeReceived.EventName,
		eventbus.StableEventID("code-received-", order.ID, fmt.Sprintf("%d", code.ReceivedAt.UnixNano())),
		order.ID,
		code.ReceivedAt,
	)
	return b.record(ctx, eventcatalog.SMSCodeReceived.Subject, &smsv1.SmsCodeReceivedEvent{
		Context: eventCtx,
		OrderId: order.ID,
		Code:    app.PublicCode(&code),
	}, eventCtx, orderAttributes(order))
}

func (b *OrderEventRecorder) OrderStatusChanged(ctx context.Context, order core.Order, previous core.OrderStatus) (eventoutbox.Record, error) {
	eventCtx := b.context(
		eventcatalog.SMSOrderStatusChanged.EventName,
		eventbus.StableEventID("order-status-", order.ID, string(previous), string(order.Status), fmt.Sprintf("%d", order.UpdatedAt.UnixNano())),
		order.ID,
		order.UpdatedAt,
	)
	return b.record(ctx, eventcatalog.SMSOrderStatusChanged.Subject, &smsv1.SmsOrderStatusChangedEvent{
		Context:        eventCtx,
		OrderId:        order.ID,
		PreviousStatus: app.PublicOrderStatus(previous),
		CurrentStatus:  app.PublicOrderStatus(order.Status),
		Error:          app.PublicError(order.LastError),
	}, eventCtx, orderAttributes(order))
}
