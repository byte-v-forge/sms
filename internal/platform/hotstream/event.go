package hotstream

import (
	"strings"
	"time"

	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
	observabilityv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/observability/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	occurredAt := cfg.OccurredAt
	if occurredAt.IsZero() {
		occurredAt = time.Now()
	}
	return &observabilityv1.HotStreamEvent{
		Metadata: &commonv1.EventMetadata{
			Id:              strings.TrimSpace(cfg.EventID),
			Type:            strings.TrimSpace(cfg.EventType),
			Version:         "v1",
			Time:            timestamppb.New(occurredAt),
			Source:          strings.TrimSpace(cfg.SourceService),
			CorrelationId:   strings.TrimSpace(cfg.CorrelationID),
			TraceId:         strings.TrimSpace(cfg.TraceID),
			SpecVersion:     "1.0",
			DataContentType: DataContentType,
		},
		ResourceType: strings.TrimSpace(cfg.ResourceType),
		ResourceId:   strings.TrimSpace(cfg.ResourceID),
		Scope:        strings.TrimSpace(cfg.Scope),
		Attributes:   CleanAttributes(cfg.Attributes),
	}
}

func ServiceStateSubject(service string) string {
	service = strings.Trim(strings.ToLower(strings.TrimSpace(service)), ".")
	if service == "" {
		service = "platform"
	}
	return SubjectPrefix + "." + service + ".state"
}

func CleanAttributes(input map[string]string) map[string]string {
	if len(input) == 0 {
		return nil
	}
	out := map[string]string{}
	for key, value := range input {
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" || value == "" {
			continue
		}
		out[key] = value
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
