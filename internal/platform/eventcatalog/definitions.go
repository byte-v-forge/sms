package eventcatalog

var (
	SMSOrderAcquired = Definition{
		Subject:      "byte.v.forge.sms.order.acquired",
		EventName:    "sms.order.acquired",
		EventVersion: EventVersionV1,
		Kind:         KindFact,
		PayloadType:  "byte.v.forge.contracts.sms.v1.SmsOrderAcquiredEvent",
		OwnerService: "sms-service",
	}
	SMSCodeReceived = Definition{
		Subject:      "byte.v.forge.sms.code.received",
		EventName:    "sms.code.received",
		EventVersion: EventVersionV1,
		Kind:         KindFact,
		PayloadType:  "byte.v.forge.contracts.sms.v1.SmsCodeReceivedEvent",
		OwnerService: "sms-service",
	}
	SMSOrderStatusChanged = Definition{
		Subject:      "byte.v.forge.sms.order.status_changed",
		EventName:    "sms.order.status_changed",
		EventVersion: EventVersionV1,
		Kind:         KindFact,
		PayloadType:  "byte.v.forge.contracts.sms.v1.SmsOrderStatusChangedEvent",
		OwnerService: "sms-service",
	}

	DeadLetter = Definition{
		Subject:      DeadLetterTopic,
		EventName:    "platform.dead_letter",
		EventVersion: EventVersionV1,
		Kind:         KindFact,
		PayloadType:  "byte.v.forge.contracts.common.v1.DeadLetterEvent",
		OwnerService: "platform",
	}
)
