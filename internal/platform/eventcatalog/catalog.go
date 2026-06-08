package eventcatalog

import commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"

const (
	StreamName      = "BYTE_V_FORGE_EVENTS"
	StreamSubject   = "byte.v.forge.>"
	DeadLetterTopic = "byte.v.forge.platform.dead_letter"
	EventVersionV1  = "v1"
)

type Kind string

const (
	KindFact    Kind = "fact"
	KindCommand Kind = "command"
)

type Definition struct {
	Subject          string
	EventName        string
	EventVersion     string
	Kind             Kind
	PayloadType      string
	OwnerService     string
	ConsumerDurable  string
	Retryable        bool
	MaxDeliveries    int
	RetryDelaySecond int
}

func (definition Definition) Proto() *commonv1.EventDefinition {
	return &commonv1.EventDefinition{
		Subject:           definition.Subject,
		EventName:         definition.EventName,
		EventVersion:      definition.EventVersion,
		Kind:              protoKind(definition.Kind),
		PayloadType:       definition.PayloadType,
		OwnerService:      definition.OwnerService,
		ConsumerDurable:   definition.ConsumerDurable,
		Retryable:         definition.Retryable,
		MaxDeliveries:     int32(definition.MaxDeliveries),
		RetryDelaySeconds: int32(definition.RetryDelaySecond),
	}
}

func Catalog() *commonv1.EventCatalog {
	definitions := All()
	out := make([]*commonv1.EventDefinition, 0, len(definitions))
	for _, definition := range definitions {
		out = append(out, definition.Proto())
	}
	return &commonv1.EventCatalog{
		StreamName:    StreamName,
		StreamSubject: StreamSubject,
		Definitions:   out,
	}
}

func protoKind(kind Kind) commonv1.EventKind {
	switch kind {
	case KindFact:
		return commonv1.EventKind_EVENT_KIND_FACT
	case KindCommand:
		return commonv1.EventKind_EVENT_KIND_COMMAND
	default:
		return commonv1.EventKind_EVENT_KIND_UNSPECIFIED
	}
}

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

func All() []Definition {
	return []Definition{
		SMSOrderAcquired,
		SMSCodeReceived,
		SMSOrderStatusChanged,
		DeadLetter,
	}
}

func Subjects() []string {
	return []string{StreamSubject}
}
