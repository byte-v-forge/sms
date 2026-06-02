package eventcatalog

import commoncatalog "github.com/byte-v-forge/common-lib/eventcatalog"

var (
	OrderAcquireRequested = commoncatalog.Definition{
		Subject:          "byte.v.forge.sms.order.acquire.requested",
		EventName:        "sms.order.acquire_requested",
		EventVersion:     commoncatalog.EventVersionV1,
		Kind:             commoncatalog.KindCommand,
		PayloadType:      "byte.v.forge.sms.internal.v1.SmsOrderAcquireRequest",
		OwnerService:     "sms-service",
		ConsumerDurable:  "sms-order-acquire",
		Retryable:        true,
		MaxDeliveries:    20,
		RetryDelaySecond: 5,
	}
	OrderPollRequested = commoncatalog.Definition{
		Subject:          "byte.v.forge.sms.order.poll.requested",
		EventName:        "sms.order.poll_requested",
		EventVersion:     commoncatalog.EventVersionV1,
		Kind:             commoncatalog.KindCommand,
		PayloadType:      "byte.v.forge.sms.internal.v1.SmsOrderPollRequest",
		OwnerService:     "sms-service",
		ConsumerDurable:  "sms-order-poll",
		Retryable:        true,
		MaxDeliveries:    20,
		RetryDelaySecond: 5,
	}
	OrderCancelRequested = commoncatalog.Definition{
		Subject:          "byte.v.forge.sms.order.cancel.requested",
		EventName:        "sms.order.cancel_requested",
		EventVersion:     commoncatalog.EventVersionV1,
		Kind:             commoncatalog.KindCommand,
		PayloadType:      "byte.v.forge.sms.internal.v1.SmsOrderCancelRequest",
		OwnerService:     "sms-service",
		ConsumerDurable:  "sms-order-cancel",
		Retryable:        true,
		MaxDeliveries:    20,
		RetryDelaySecond: 5,
	}
)
