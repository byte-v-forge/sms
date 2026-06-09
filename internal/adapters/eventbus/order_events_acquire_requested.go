package eventbusadapter

import (
	"context"
	"strings"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/app"
	"github.com/byte-v-forge/sms/internal/core"
	smseventcatalog "github.com/byte-v-forge/sms/internal/eventcatalog"
	"github.com/byte-v-forge/sms/internal/platform/eventbus"
	"github.com/byte-v-forge/sms/internal/platform/eventoutbox"
)

func (b *OrderEventRecorder) recordOrderAcquireRequested(ctx context.Context, order core.Order, route core.Route, reason string) (eventoutbox.Record, error) {
	reason = strings.TrimSpace(reason)
	metadata := b.metadata(
		smseventcatalog.OrderAcquireRequested.EventName,
		smseventcatalog.OrderAcquireRequested.Subject,
		eventbus.StableEventID("order-acquire-", order.ID, order.RequestID, reason),
		order.ID,
		order.UpdatedAt,
	)
	return b.record(ctx, smseventcatalog.OrderAcquireRequested, &smsinternalv1.SmsOrderAcquireRequest{
		OrderId:       order.ID,
		RequestId:     order.RequestID,
		Reason:        reason,
		AcquireParams: app.PublicAcquireParamsFromRoute(route),
	}, metadata, eventbus.WithNonEmptyAttribute(orderAttributes(order), "reason", reason))
}
