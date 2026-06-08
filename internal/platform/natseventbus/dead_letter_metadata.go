package natseventbus

import (
	"fmt"

	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
	"github.com/byte-v-forge/sms/internal/platform/eventbus"
	"github.com/byte-v-forge/sms/internal/platform/eventcatalog"
)

func deadLetterMetadata(envelope *commonv1.EventEnvelope, original originalEventInfo, durable string, attempt int32) *commonv1.EventMetadata {
	eventID := eventbus.StableEventID("dead-letter-", envelope.GetSubject(), original.id, durable, fmt.Sprintf("%d", attempt))
	return eventbus.NewEventMetadata(eventbus.EventMetadataConfig{
		EventID:       eventID,
		EventName:     "platform.dead_letter",
		EventVersion:  eventcatalog.EventVersionV1,
		SourceService: "platform-eventbus",
		Subject:       eventcatalog.DeadLetter.Subject,
		CorrelationID: original.correlationID,
		TraceID:       original.traceID,
	})
}
