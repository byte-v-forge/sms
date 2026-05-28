package eventbusadapter

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/byte-v-forge/common-lib/eventbus"
	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/app"
	"github.com/byte-v-forge/sms/internal/core"
)

const defaultCancelRetryDelay = 30 * time.Second

type OrderCancelWorker struct {
	service *app.OrderService
}

func RunOrderCancelWorker(ctx context.Context, consumer eventbus.Consumer, service *app.OrderService) error {
	worker := &OrderCancelWorker{service: service}
	return eventbus.RunConsumerWorker(ctx, eventbus.ConsumerWorkerConfig{
		Name:     "sms order cancel requests",
		Consumer: consumer,
		Handler:  worker.handle,
	})
}

func (w *OrderCancelWorker) handle(ctx context.Context, message eventbus.ReceivedMessage) {
	request, ok := decodeCancelRequest(message)
	if !ok {
		eventbus.TermMessage(ctx, message, "terminate malformed sms order cancel request", nil)
		return
	}
	order, err := w.service.RunCancelRequest(ctx, request.GetOrderId(), request.GetRequestId())
	if err == nil {
		eventbus.AckMessage(ctx, message, "ack sms order cancel request", nil)
		return
	}
	log.Printf("sms order cancel failed: order_id=%s error=%v", request.GetOrderId(), err)
	if delay, ok := cancelRetryDelay(err, time.Now()); ok {
		eventbus.NakMessageDelay(ctx, message, delay, "delay sms order cancel retry", nil)
		return
	}
	if cancelErrorRetryable(err) {
		delay := w.cancelDelay(order)
		eventbus.NakMessageDelay(ctx, message, delay, "retry sms order cancel", nil)
		return
	}
	eventbus.TermMessage(ctx, message, "terminate non-retryable sms order cancel request", nil)
}

func decodeCancelRequest(message eventbus.ReceivedMessage) (*smsinternalv1.SmsOrderCancelRequest, bool) {
	request := &smsinternalv1.SmsOrderCancelRequest{}
	if err := eventbus.UnmarshalPayload(message, request); err != nil {
		log.Printf("decode sms order cancel request failed: event_id=%s error=%v", eventbus.EventID(message), err)
		return nil, false
	}
	if strings.TrimSpace(request.GetOrderId()) == "" {
		log.Printf("sms order cancel request missing order_id: event_id=%s", eventbus.EventID(message))
		return nil, false
	}
	return request, true
}

func cancelRetryDelay(err error, now time.Time) (time.Duration, bool) {
	var retryErr *app.CancelRetryError
	if !errors.As(err, &retryErr) {
		return 0, false
	}
	delay := retryErr.RetryAt.Sub(now)
	if retryErr.RetryAt.IsZero() || delay <= 0 {
		delay = defaultCancelRetryDelay
	}
	return delay, true
}

func cancelErrorRetryable(err error) bool {
	var smsErr *core.Error
	if !errors.As(err, &smsErr) {
		return true
	}
	switch smsErr.Code {
	case core.CodeOrderNotFound, core.CodeOrderAlreadyFinalized, core.CodeOrderExpired, core.CodeCancelNotAllowed:
		return false
	default:
		return smsErr.Retryable
	}
}

func (w *OrderCancelWorker) cancelDelay(order core.Order) time.Duration {
	if order.ID != "" && !order.CancelAllowedAt.IsZero() {
		delay := time.Until(order.CancelAllowedAt)
		if delay > 0 {
			return delay
		}
	}
	return defaultCancelRetryDelay
}
