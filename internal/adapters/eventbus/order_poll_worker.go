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
)

const defaultPollRetryDelay = 5 * time.Second

type OrderPollWorker struct {
	service *app.OrderService
}

func RunOrderPollWorker(ctx context.Context, consumer eventbus.Consumer, service *app.OrderService) error {
	worker := &OrderPollWorker{service: service}
	return eventbus.RunTypedConsumerWorker(ctx, eventbus.TypedConsumerWorkerConfig[*smsinternalv1.SmsOrderPollRequest]{
		Name:           "sms order poll requests",
		Consumer:       consumer,
		NewMessage:     func() *smsinternalv1.SmsOrderPollRequest { return &smsinternalv1.SmsOrderPollRequest{} },
		Validate:       func(request *smsinternalv1.SmsOrderPollRequest) error { return validateOrderID(request.GetOrderId()) },
		Handler:        worker.handle,
		MalformedLabel: "terminate malformed sms order poll request",
	})
}

func (w *OrderPollWorker) handle(ctx context.Context, request *smsinternalv1.SmsOrderPollRequest, _ eventbus.ReceivedMessage) eventbus.HandlerResult {
	order, code, err := w.service.CheckCode(ctx, request.GetOrderId())
	if err != nil {
		log.Printf("sms order poll failed: order_id=%s error=%v", request.GetOrderId(), err)
		if pollErrorRetryable(err) {
			delay := w.pollDelay(ctx, order)
			return eventbus.NakResult(delay, "delay sms order poll retry")
		}
		return eventbus.TermResult("terminate non-retryable sms order poll request")
	}
	if code != nil || order.Status == core.StatusCodeReceived || order.Status.IsFinal() {
		return eventbus.AckResult("ack sms order poll request")
	}
	return eventbus.NakResult(w.pollDelay(ctx, order), "delay sms order poll")
}

func (w *OrderPollWorker) pollDelay(ctx context.Context, order core.Order) time.Duration {
	if order.ID == "" {
		return defaultPollRetryDelay
	}
	delay := w.service.PollInterval(ctx, order)
	if delay <= 0 {
		return defaultPollRetryDelay
	}
	return delay
}

func pollErrorRetryable(err error) bool {
	var smsErr *core.Error
	if !errors.As(err, &smsErr) {
		return true
	}
	switch smsErr.Code {
	case core.CodeOrderNotFound, core.CodeOrderAlreadyFinalized, core.CodeOrderExpired:
		return false
	default:
		return smsErr.Retryable
	}
}
