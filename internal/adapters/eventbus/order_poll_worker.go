package eventbusadapter

import (
	"context"
	"time"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/app"
	smseventcatalog "github.com/byte-v-forge/sms/internal/eventcatalog"
	"github.com/byte-v-forge/sms/internal/platform/eventbus"
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
		Expected:       smseventcatalog.OrderPollRequested.ExpectedMessage(),
		NewMessage:     func() *smsinternalv1.SmsOrderPollRequest { return &smsinternalv1.SmsOrderPollRequest{} },
		Validate:       func(request *smsinternalv1.SmsOrderPollRequest) error { return validateOrderID(request.GetOrderId()) },
		Handler:        worker.handle,
		MalformedLabel: "terminate malformed sms order poll request",
	})
}
