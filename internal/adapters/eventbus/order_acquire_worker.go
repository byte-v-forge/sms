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

const defaultAcquireRetryDelay = 15 * time.Second

type OrderAcquireWorker struct {
	service *app.OrderService
}

func RunOrderAcquireWorker(ctx context.Context, consumer eventbus.Consumer, service *app.OrderService) error {
	worker := &OrderAcquireWorker{service: service}
	return eventbus.RunConsumerWorker(ctx, eventbus.ConsumerWorkerConfig{
		Name:     "sms order acquire requests",
		Consumer: consumer,
		Handler:  worker.handle,
	})
}

func (w *OrderAcquireWorker) handle(ctx context.Context, message eventbus.ReceivedMessage) {
	request, ok := decodeAcquireRequest(message)
	if !ok {
		eventbus.TermMessage(ctx, message, "terminate malformed sms order acquire request", nil)
		return
	}
	_, err := w.service.RunAcquireRequest(ctx, request.GetOrderId(), request.GetRequestId(), app.RouteFromPublicAcquireParams(request.GetAcquireParams()))
	if err == nil {
		eventbus.AckMessage(ctx, message, "ack sms order acquire request", nil)
		return
	}
	log.Printf("sms order acquire failed: order_id=%s error=%v", request.GetOrderId(), err)
	if acquireErrorRetryable(err) {
		eventbus.NakMessageDelay(ctx, message, defaultAcquireRetryDelay, "retry sms order acquire", nil)
		return
	}
	eventbus.TermMessage(ctx, message, "terminate non-retryable sms order acquire request", nil)
}

func decodeAcquireRequest(message eventbus.ReceivedMessage) (*smsinternalv1.SmsOrderAcquireRequest, bool) {
	request := &smsinternalv1.SmsOrderAcquireRequest{}
	if err := eventbus.UnmarshalPayload(message, request); err != nil {
		log.Printf("decode sms order acquire request failed: event_id=%s error=%v", eventbus.EventID(message), err)
		return nil, false
	}
	if strings.TrimSpace(request.GetOrderId()) == "" {
		log.Printf("sms order acquire request missing order_id: event_id=%s", eventbus.EventID(message))
		return nil, false
	}
	return request, true
}

func acquireErrorRetryable(err error) bool {
	var smsErr *core.Error
	if !errors.As(err, &smsErr) {
		return true
	}
	switch smsErr.Code {
	case core.CodeValidationFailed, core.CodeUnsupportedOperation, core.CodeInsufficientBalance, core.CodeOrderExpired, core.CodeOrderAlreadyFinalized:
		return false
	default:
		return smsErr.Retryable
	}
}
