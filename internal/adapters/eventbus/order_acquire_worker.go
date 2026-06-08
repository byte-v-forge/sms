package eventbusadapter

import (
	"context"
	"log"
	"time"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/app"
	"github.com/byte-v-forge/sms/internal/core"
	smseventcatalog "github.com/byte-v-forge/sms/internal/eventcatalog"
	"github.com/byte-v-forge/sms/internal/platform/eventbus"
)

const defaultAcquireRetryDelay = 15 * time.Second

type OrderAcquireWorker struct {
	service *app.OrderService
}

func RunOrderAcquireWorker(ctx context.Context, consumer eventbus.Consumer, service *app.OrderService) error {
	worker := &OrderAcquireWorker{service: service}
	return eventbus.RunTypedConsumerWorker(ctx, eventbus.TypedConsumerWorkerConfig[*smsinternalv1.SmsOrderAcquireRequest]{
		Name:       "sms order acquire requests",
		Consumer:   consumer,
		Expected:   smseventcatalog.OrderAcquireRequested.ExpectedMessage(),
		NewMessage: func() *smsinternalv1.SmsOrderAcquireRequest { return &smsinternalv1.SmsOrderAcquireRequest{} },
		Validate: func(request *smsinternalv1.SmsOrderAcquireRequest) error {
			return validateOrderID(request.GetOrderId())
		},
		Handler:        worker.handle,
		MalformedLabel: "terminate malformed sms order acquire request",
	})
}

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

func acquireErrorRetryable(err error) bool {
	return coreErrorRetryableUnless(
		err,
		core.CodeValidationFailed,
		core.CodeUnsupportedOperation,
		core.CodeInsufficientBalance,
		core.CodeOrderExpired,
		core.CodeOrderAlreadyFinalized,
	)
}
