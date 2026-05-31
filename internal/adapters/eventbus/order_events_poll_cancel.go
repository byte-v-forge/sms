package eventbusadapter

import (
	"context"
	"fmt"
	"strings"

	"github.com/byte-v-forge/common-lib/eventbus"
	"github.com/byte-v-forge/common-lib/eventcatalog"
	"github.com/byte-v-forge/common-lib/eventoutbox"
	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

func (b *OrderEventRecorder) OrderPollRequested(ctx context.Context, order core.Order, reason string) (eventoutbox.Record, error) {
	reason = strings.TrimSpace(reason)
	eventCtx := b.context(
		eventcatalog.SMSOrderPollRequested.EventName,
		eventbus.StableEventID("order-poll-", order.ID, reason, fmt.Sprintf("%d", order.UpdatedAt.UnixNano())),
		order.ID,
		order.UpdatedAt,
	)
	return b.record(ctx, eventcatalog.SMSOrderPollRequested.Subject, &smsinternalv1.SmsOrderPollRequest{
		OrderId: order.ID,
		Reason:  reason,
	}, eventCtx, eventbus.WithNonEmptyAttribute(orderAttributes(order), "reason", reason))
}

func (b *OrderEventRecorder) OrderCancelRequested(ctx context.Context, order core.Order, requestID string, reason string) (eventoutbox.Record, error) {
	requestID = strings.TrimSpace(requestID)
	reason = strings.TrimSpace(reason)
	eventCtx := b.context(
		eventcatalog.SMSOrderCancelRequested.EventName,
		eventbus.StableEventID("order-cancel-", order.ID, requestID, reason),
		order.ID,
		order.UpdatedAt,
	)
	return b.record(ctx, eventcatalog.SMSOrderCancelRequested.Subject, &smsinternalv1.SmsOrderCancelRequest{
		OrderId:   order.ID,
		RequestId: requestID,
		Reason:    reason,
	}, eventCtx, eventbus.WithNonEmptyAttribute(eventbus.WithNonEmptyAttribute(orderAttributes(order), "request_id", requestID), "reason", reason))
}
