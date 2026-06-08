package natseventbus

import (
	"fmt"

	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
	"github.com/byte-v-forge/sms/internal/platform/eventbus"
	"github.com/byte-v-forge/sms/internal/platform/eventcatalog"
)

func deadLetterMessage(durable string, envelope *commonv1.EventEnvelope, attempt int32, reason string) (eventbus.Message, error) {
	original := originalEvent(envelope)
	metadata := deadLetterMetadata(envelope, original, durable, attempt)
	return eventcatalog.DeadLetter.NewMessage(
		&commonv1.DeadLetterEvent{
			Metadata:             metadata,
			OriginalSubject:      envelope.GetSubject(),
			OriginalEventId:      original.id,
			OriginalEventType:    original.eventName,
			OriginalEventVersion: original.eventVersion,
			OriginalSource:       original.source,
			ConsumerDurable:      durable,
			DeliveryAttempt:      attempt,
			ErrorCode:            "terminated",
			ErrorMessage:         reason,
			CorrelationId:        original.correlationID,
		},
		metadata,
		eventbus.Attributes(
			"original_subject", envelope.GetSubject(),
			"original_event_id", original.id,
			"consumer_durable", durable,
			"delivery_attempt", fmt.Sprintf("%d", attempt),
		),
	)
}
