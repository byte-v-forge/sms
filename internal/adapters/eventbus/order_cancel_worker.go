package eventbusadapter

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/byte-v-forge/common-lib/eventbus"
	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/app"
	"github.com/byte-v-forge/sms/internal/core"
	smseventcatalog "github.com/byte-v-forge/sms/internal/eventcatalog"
)

const defaultCancelRetryDelay = 30 * time.Second

type OrderCancelWorker struct {
	service *app.OrderService
}

func RunOrderCancelWorker(ctx context.Context, consumer eventbus.Consumer, service *app.OrderService) error {
	worker := &OrderCancelWorker{service: service}
	return eventbus.RunTypedConsumerWorker(ctx, eventbus.TypedConsumerWorkerConfig[*smsinternalv1.SmsOrderCancelRequest]{
		Name:           "sms order cancel requests",
		Consumer:       consumer,
		Expected:       smseventcatalog.OrderCancelRequested.ExpectedMessage(),
		NewMessage:     func() *smsinternalv1.SmsOrderCancelRequest { return &smsinternalv1.SmsOrderCancelRequest{} },
		Validate:       func(request *smsinternalv1.SmsOrderCancelRequest) error { return validateOrderID(request.GetOrderId()) },
		Handler:        worker.handle,
		MalformedLabel: "terminate malformed sms order cancel request",
	})
}

func (w *OrderCancelWorker) handle(ctx context.Context, request *smsinternalv1.SmsOrderCancelRequest, _ eventbus.ReceivedMessage) eventbus.HandlerResult {
	order, err := w.service.RunCancelRequest(ctx, request.GetOrderId(), request.GetRequestId())
	if err == nil {
		return eventbus.AckResult("ack sms order cancel request")
	}
	log.Printf("sms order cancel failed: order_id=%s error=%v", request.GetOrderId(), err)
	if delay, ok := cancelRetryDelay(err, time.Now()); ok {
		return eventbus.NakResult(delay, "delay sms order cancel retry")
	}
	if cancelErrorRetryable(err) {
		return eventbus.NakResult(w.cancelDelay(order), "retry sms order cancel")
	}
	return eventbus.TermResult("terminate non-retryable sms order cancel request")
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
