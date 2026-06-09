package eventbusadapter

import (
	"context"
	"time"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/app"
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
		Name:           "sms order acquire requests",
		Consumer:       consumer,
		Expected:       smseventcatalog.OrderAcquireRequested.ExpectedMessage(),
		NewMessage:     func() *smsinternalv1.SmsOrderAcquireRequest { return &smsinternalv1.SmsOrderAcquireRequest{} },
		Validate:       validateAcquireRequest,
		Handler:        worker.handle,
		MalformedLabel: "terminate malformed sms order acquire request",
	})
}
