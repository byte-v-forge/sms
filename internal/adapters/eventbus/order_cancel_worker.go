package eventbusadapter

import (
	"context"
	"time"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/app"
	smseventcatalog "github.com/byte-v-forge/sms/internal/eventcatalog"
	"github.com/byte-v-forge/sms/internal/platform/eventbus"
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
