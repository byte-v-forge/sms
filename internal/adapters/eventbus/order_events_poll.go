package eventbusadapter

import (
	"context"
	"strings"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
	smseventcatalog "github.com/byte-v-forge/sms/internal/eventcatalog"
	"github.com/byte-v-forge/sms/internal/platform/eventbus"
	"github.com/byte-v-forge/sms/internal/platform/eventoutbox"
)

func (b *OrderEventRecorder) recordOrderPollRequested(ctx context.Context, order core.Order, reason string) (eventoutbox.Record, error) {
	reason = strings.TrimSpace(reason)
	metadata := b.metadata(
		smseventcatalog.OrderPollRequested.EventName,
		smseventcatalog.OrderPollRequested.Subject,
		eventbus.StableEventID("order-poll-", order.ID, reason, eventTimeSuffix(order.UpdatedAt)),
		order.ID,
		order.UpdatedAt,
	)
	return b.record(ctx, smseventcatalog.OrderPollRequested, &smsinternalv1.SmsOrderPollRequest{
		OrderId: order.ID,
		Reason:  reason,
	}, metadata, eventbus.WithNonEmptyAttribute(orderAttributes(order), "reason", reason))
}
