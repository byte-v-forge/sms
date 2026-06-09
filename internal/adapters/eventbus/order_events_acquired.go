package eventbusadapter

import (
	"context"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/app"
	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/eventbus"
	"github.com/byte-v-forge/sms/internal/platform/eventcatalog"
	"github.com/byte-v-forge/sms/internal/platform/eventoutbox"
)

func (b *OrderEventRecorder) recordOrderAcquired(ctx context.Context, order core.Order) (eventoutbox.Record, error) {
	metadata := b.metadata(
		eventcatalog.SMSOrderAcquired.EventName,
		eventcatalog.SMSOrderAcquired.Subject,
		eventbus.StableEventID("order-acquired-", order.ID, order.UpstreamOrderID),
		order.ID,
		order.AcquiredAt,
	)
	return b.record(ctx, eventcatalog.SMSOrderAcquired, &smsv1.SmsOrderAcquiredEvent{
		Metadata: metadata,
		Order:    app.PublicOrder(order),
	}, metadata, orderAttributes(order))
}
