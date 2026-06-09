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

func (b *OrderEventRecorder) recordOrderCancelRequested(ctx context.Context, order core.Order, requestID string, reason string) (eventoutbox.Record, error) {
	requestID = strings.TrimSpace(requestID)
	reason = strings.TrimSpace(reason)
	metadata := b.metadata(
		smseventcatalog.OrderCancelRequested.EventName,
		smseventcatalog.OrderCancelRequested.Subject,
		eventbus.StableEventID("order-cancel-", order.ID, requestID, reason),
		order.ID,
		order.UpdatedAt,
	)
	attrs := eventbus.WithNonEmptyAttribute(orderAttributes(order), "request_id", requestID)
	attrs = eventbus.WithNonEmptyAttribute(attrs, "reason", reason)
	return b.record(ctx, smseventcatalog.OrderCancelRequested, &smsinternalv1.SmsOrderCancelRequest{
		OrderId:   order.ID,
		RequestId: requestID,
		Reason:    reason,
	}, metadata, attrs)
}
