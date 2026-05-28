package eventbusadapter

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/byte-v-forge/common-lib/eventbus"
	"github.com/byte-v-forge/common-lib/eventcatalog"
	"github.com/byte-v-forge/common-lib/eventoutbox"
	commonv1 "github.com/byte-v-forge/common-lib/gen/go/byte/v/forge/contracts/common/v1"
	smsv1 "github.com/byte-v-forge/common-lib/gen/go/byte/v/forge/contracts/sms/v1"
	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/app"
	"github.com/byte-v-forge/sms/internal/core"
	"google.golang.org/protobuf/proto"
)

const defaultSourceService = "sms-service"

type OrderEventRecorder struct {
	source string
}

func NewOrderEventRecorder(source string) *OrderEventRecorder {
	source = strings.TrimSpace(source)
	if source == "" {
		source = defaultSourceService
	}
	return &OrderEventRecorder{source: source}
}

func (b *OrderEventRecorder) OrderAcquireRequested(ctx context.Context, order core.Order, route core.Route, reason string) (eventoutbox.Record, error) {
	reason = strings.TrimSpace(reason)
	eventCtx := b.context(
		eventcatalog.SMSOrderAcquireRequested.EventName,
		eventbus.StableEventID("order-acquire-", order.ID, order.RequestID, reason),
		order.ID,
		order.UpdatedAt,
	)
	return b.record(ctx, eventcatalog.SMSOrderAcquireRequested.Subject, &smsinternalv1.SmsOrderAcquireRequest{
		OrderId:   order.ID,
		RequestId:      order.RequestID,
		Reason:         reason,
		AcquireParams:  app.PublicAcquireParamsFromRoute(route),
	}, eventCtx, eventbus.WithNonEmptyAttribute(orderAttributes(order), "reason", reason))
}

func (b *OrderEventRecorder) OrderAcquired(ctx context.Context, order core.Order) (eventoutbox.Record, error) {
	eventCtx := b.context(
		eventcatalog.SMSOrderAcquired.EventName,
		eventbus.StableEventID("order-acquired-", order.ID, order.UpstreamOrderID),
		order.ID,
		order.AcquiredAt,
	)
	return b.record(ctx, eventcatalog.SMSOrderAcquired.Subject, &smsv1.SmsOrderAcquiredEvent{
		Context:    eventCtx,
		Order: app.PublicOrder(order),
	}, eventCtx, orderAttributes(order))
}

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
		Reason:       reason,
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
		OrderId: order.ID,
		RequestId:    requestID,
		Reason:       reason,
	}, eventCtx, eventbus.WithNonEmptyAttribute(eventbus.WithNonEmptyAttribute(orderAttributes(order), "request_id", requestID), "reason", reason))
}

func (b *OrderEventRecorder) CodeReceived(ctx context.Context, order core.Order, code core.SMSCode) (eventoutbox.Record, error) {
	eventCtx := b.context(
		eventcatalog.SMSCodeReceived.EventName,
		eventbus.StableEventID("code-received-", order.ID, code.Value, fmt.Sprintf("%d", code.ReceivedAt.UnixNano())),
		order.ID,
		code.ReceivedAt,
	)
	return b.record(ctx, eventcatalog.SMSCodeReceived.Subject, &smsv1.SmsCodeReceivedEvent{
		Context:      eventCtx,
		OrderId: order.ID,
		Code:         app.PublicCode(&code),
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
		OrderId:   order.ID,
		PreviousStatus: app.PublicOrderStatus(previous),
		CurrentStatus:  app.PublicOrderStatus(order.Status),
		Error:          app.PublicError(order.LastError),
	}, eventCtx, orderAttributes(order))
}

func (b *OrderEventRecorder) record(_ context.Context, subject string, message proto.Message, eventCtx *commonv1.EventContext, attrs map[string]string) (eventoutbox.Record, error) {
	if b == nil {
		return eventoutbox.Record{}, nil
	}
	return eventoutbox.NewRecord(eventbus.Message{
		Subject:    subject,
		Event:      message,
		Context:    eventCtx,
		Attributes: attrs,
	})
}

func (b *OrderEventRecorder) context(eventName string, eventID string, correlationID string, occurredAt time.Time) *commonv1.EventContext {
	return eventbus.NewEventContext(eventbus.EventContextConfig{
		EventID:       eventID,
		EventName:     eventName,
		EventVersion:  eventcatalog.EventVersionV1,
		OccurredAt:    occurredAt,
		SourceService: b.source,
		CorrelationID: correlationID,
	})
}

func orderAttributes(order core.Order) map[string]string {
	return eventbus.Attributes(
		"order_id", order.ID,
		"provider_key", order.ProviderKey,
		"status", string(order.Status),
	)
}
