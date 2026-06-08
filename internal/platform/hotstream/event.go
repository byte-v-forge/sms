package hotstream

import (
	"time"

	observabilityv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/observability/v1"
)

const SubjectPrefix = "byte.v.forge.hot"
const DataContentType = "application/x-protobuf"

type EventConfig struct {
	EventID       string
	EventType     string
	SourceService string
	ResourceType  string
	ResourceID    string
	Scope         string
	OccurredAt    time.Time
	CorrelationID string
	TraceID       string
	Attributes    map[string]string
}

func NewEvent(cfg EventConfig) *observabilityv1.HotStreamEvent {
	return &observabilityv1.HotStreamEvent{
		Metadata:     eventMetadata(cfg),
		ResourceType: cleanEventText(cfg.ResourceType),
		ResourceId:   cleanEventText(cfg.ResourceID),
		Scope:        cleanEventText(cfg.Scope),
		Attributes:   CleanAttributes(cfg.Attributes),
	}
}
