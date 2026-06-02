package eventbusadapter

import (
	"context"
	"strings"

	"github.com/byte-v-forge/common-lib/eventbus"
	"github.com/byte-v-forge/common-lib/eventcatalog"
	"github.com/byte-v-forge/common-lib/eventoutbox"
	smsv1 "github.com/byte-v-forge/common-lib/gen/go/byte/v/forge/contracts/sms/v1"
	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/app"
	"github.com/byte-v-forge/sms/internal/core"
	smseventcatalog "github.com/byte-v-forge/sms/internal/eventcatalog"
)

func (b *OrderEventRecorder) OrderAcquireRequested(ctx context.Context, order core.Order, route core.Route, reason string) (eventoutbox.Record, error) {
	reason = strings.TrimSpace(reason)
	eventCtx := b.context(
		smseventcatalog.OrderAcquireRequested.EventName,
		eventbus.StableEventID("order-acquire-", order.ID, order.RequestID, reason),
		order.ID,
		order.UpdatedAt,
	)
	return b.record(ctx, smseventcatalog.OrderAcquireRequested, &smsinternalv1.SmsOrderAcquireRequest{
		OrderId:       order.ID,
		RequestId:     order.RequestID,
		Reason:        reason,
		AcquireParams: app.PublicAcquireParamsFromRoute(route),
	}, eventCtx, eventbus.WithNonEmptyAttribute(orderAttributes(order), "reason", reason))
}

func (b *OrderEventRecorder) OrderAcquired(ctx context.Context, order core.Order) (eventoutbox.Record, error) {
	eventCtx := b.context(
		eventcatalog.SMSOrderAcquired.EventName,
		eventbus.StableEventID("order-acquired-", order.ID, order.UpstreamOrderID),
		order.ID,
		order.AcquiredAt,
	)
	return b.record(ctx, eventcatalog.SMSOrderAcquired, &smsv1.SmsOrderAcquiredEvent{
		Context: eventCtx,
		Order:   app.PublicOrder(order),
	}, eventCtx, orderAttributes(order))
}
