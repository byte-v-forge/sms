package eventbusadapter

import (
	"context"
	"fmt"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/app"
	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/eventbus"
	"github.com/byte-v-forge/sms/internal/platform/eventcatalog"
	"github.com/byte-v-forge/sms/internal/platform/eventoutbox"
)

func (b *OrderEventRecorder) CodeReceived(ctx context.Context, order core.Order, code core.SMSCode) (eventoutbox.Record, error) {
	metadata := b.metadata(
		eventcatalog.SMSCodeReceived.EventName,
		eventcatalog.SMSCodeReceived.Subject,
		eventbus.StableEventID("code-received-", order.ID, fmt.Sprintf("%d", code.ReceivedAt.UnixNano())),
		order.ID,
		code.ReceivedAt,
	)
	return b.record(ctx, eventcatalog.SMSCodeReceived, &smsv1.SmsCodeReceivedEvent{
		Metadata: metadata,
		OrderId:  order.ID,
		Code:     app.PublicCode(&code),
	}, metadata, orderAttributes(order))
}

func (b *OrderEventRecorder) OrderStatusChanged(ctx context.Context, order core.Order, previous core.OrderStatus) (eventoutbox.Record, error) {
	metadata := b.metadata(
		eventcatalog.SMSOrderStatusChanged.EventName,
		eventcatalog.SMSOrderStatusChanged.Subject,
		eventbus.StableEventID("order-status-", order.ID, string(previous), string(order.Status), fmt.Sprintf("%d", order.UpdatedAt.UnixNano())),
		order.ID,
		order.UpdatedAt,
	)
	return b.record(ctx, eventcatalog.SMSOrderStatusChanged, &smsv1.SmsOrderStatusChangedEvent{
		Metadata:       metadata,
		OrderId:        order.ID,
		PreviousStatus: app.PublicOrderStatus(previous),
		CurrentStatus:  app.PublicOrderStatus(order.Status),
		Error:          app.PublicError(order.LastError),
	}, metadata, orderAttributes(order))
}
