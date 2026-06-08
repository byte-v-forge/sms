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
