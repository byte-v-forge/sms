package eventbusadapter

import (
	"context"
	"log"
	"time"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/platform/eventbus"
)

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
