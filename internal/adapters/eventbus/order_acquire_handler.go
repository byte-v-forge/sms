package eventbusadapter

import (
	"context"
	"log"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/app"
	"github.com/byte-v-forge/sms/internal/platform/eventbus"
)

func (w *OrderAcquireWorker) handle(ctx context.Context, request *smsinternalv1.SmsOrderAcquireRequest, _ eventbus.ReceivedMessage) eventbus.HandlerResult {
	_, err := w.service.RunAcquireRequest(ctx, request.GetOrderId(), request.GetRequestId(), app.RouteFromPublicAcquireParams(request.GetAcquireParams()))
	if err == nil {
		return eventbus.AckResult("ack sms order acquire request")
	}
	log.Printf("sms order acquire failed: order_id=%s error=%v", request.GetOrderId(), err)
	if acquireErrorRetryable(err) {
		return eventbus.NakResult(defaultAcquireRetryDelay, "retry sms order acquire")
	}
	return eventbus.TermResult("terminate non-retryable sms order acquire request")
}
