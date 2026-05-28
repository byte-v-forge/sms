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

const defaultPollRetryDelay = 5 * time.Second

type OrderPollWorker struct {
	service *app.OrderService
}

func RunOrderPollWorker(ctx context.Context, consumer eventbus.Consumer, service *app.OrderService) error {
	worker := &OrderPollWorker{service: service}
	return eventbus.RunConsumerWorker(ctx, eventbus.ConsumerWorkerConfig{
		Name:     "sms order poll requests",
		Consumer: consumer,
		Handler:  worker.handle,
	})
}

func (w *OrderPollWorker) handle(ctx context.Context, message eventbus.ReceivedMessage) {
	request, ok := decodePollRequest(message)
	if !ok {
		eventbus.TermMessage(ctx, message, "terminate malformed sms order poll request", nil)
		return
	}
	order, code, err := w.service.CheckCode(ctx, request.GetOrderId())
	if err != nil {
		log.Printf("sms order poll failed: order_id=%s error=%v", request.GetOrderId(), err)
		if pollErrorRetryable(err) {
			delay := w.pollDelay(ctx, order)
			eventbus.NakMessageDelay(ctx, message, delay, "delay sms order poll retry", nil)
			return
		}
		eventbus.TermMessage(ctx, message, "terminate non-retryable sms order poll request", nil)
		return
	}
	if code != nil || order.Status == core.StatusCodeReceived || order.Status.IsFinal() {
		eventbus.AckMessage(ctx, message, "ack sms order poll request", nil)
		return
	}
	eventbus.NakMessageDelay(ctx, message, w.pollDelay(ctx, order), "delay sms order poll", nil)
}

func decodePollRequest(message eventbus.ReceivedMessage) (*smsinternalv1.SmsOrderPollRequest, bool) {
	request := &smsinternalv1.SmsOrderPollRequest{}
	if err := eventbus.UnmarshalPayload(message, request); err != nil {
		log.Printf("decode sms order poll request failed: event_id=%s error=%v", eventbus.EventID(message), err)
		return nil, false
	}
	if strings.TrimSpace(request.GetOrderId()) == "" {
		log.Printf("sms order poll request missing order_id: event_id=%s", eventbus.EventID(message))
		return nil, false
	}
	return request, true
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
