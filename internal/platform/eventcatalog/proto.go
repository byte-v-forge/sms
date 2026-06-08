package eventcatalog

import commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"

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
