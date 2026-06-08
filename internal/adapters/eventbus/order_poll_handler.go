package eventbusadapter

import (
	"context"
	"log"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/eventbus"
)

func (w *OrderPollWorker) handle(ctx context.Context, request *smsinternalv1.SmsOrderPollRequest, _ eventbus.ReceivedMessage) eventbus.HandlerResult {
	order, code, err := w.service.CheckCode(ctx, request.GetOrderId())
	if err != nil {
		log.Printf("sms order poll failed: order_id=%s error=%v", request.GetOrderId(), err)
		if pollErrorRetryable(err) {
			return eventbus.NakResult(w.pollDelay(ctx, order), "delay sms order poll retry")
		}
		return eventbus.TermResult("terminate non-retryable sms order poll request")
	}
	if code != nil || order.Status == core.StatusCodeReceived || order.Status.IsFinal() {
		return eventbus.AckResult("ack sms order poll request")
	}
	return eventbus.NakResult(w.pollDelay(ctx, order), "delay sms order poll")
}
