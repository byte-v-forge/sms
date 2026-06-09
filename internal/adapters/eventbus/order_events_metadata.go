package eventbusadapter

import (
	"time"

	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
	"github.com/byte-v-forge/sms/internal/platform/eventbus"
	"github.com/byte-v-forge/sms/internal/platform/eventcatalog"
)

func newOrderEventMetadata(source string, eventName string, subject string, eventID string, correlationID string, occurredAt time.Time) *commonv1.EventMetadata {
	return eventbus.NewEventMetadata(eventbus.EventMetadataConfig{
		EventID:       eventID,
		EventName:     eventName,
		EventVersion:  eventcatalog.EventVersionV1,
		OccurredAt:    occurredAt,
		SourceService: source,
		Subject:       subject,
		CorrelationID: correlationID,
	})
}
