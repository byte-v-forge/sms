package eventbus

import (
	"strings"
	"time"

	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
	"github.com/byte-v-forge/sms/internal/platform/hashx"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const DefaultEventVersion = "v1"
const DefaultEventSpecVersion = "1.0"

type EventMetadataConfig struct {
	EventID        string
	EventName      string
	EventVersion   string
	OccurredAt     time.Time
	SourceService  string
	Subject        string
	CorrelationID  string
	TraceID        string
	IdempotencyKey string
	DataSchema     string
}

func NewEventMetadata(cfg EventMetadataConfig) *commonv1.EventMetadata {
	eventVersion := strings.TrimSpace(cfg.EventVersion)
	if eventVersion == "" {
		eventVersion = DefaultEventVersion
	}
	occurredAt := cfg.OccurredAt
	if occurredAt.IsZero() {
		occurredAt = time.Now()
	}
	idempotencyKey := strings.TrimSpace(cfg.IdempotencyKey)
	eventID := strings.TrimSpace(cfg.EventID)
	if idempotencyKey == "" {
		idempotencyKey = eventID
	}
	return &commonv1.EventMetadata{
		Id:              eventID,
		Type:            strings.TrimSpace(cfg.EventName),
		Version:         eventVersion,
		Time:            timestamppb.New(occurredAt),
		Source:          strings.TrimSpace(cfg.SourceService),
		CorrelationId:   strings.TrimSpace(cfg.CorrelationID),
		TraceId:         strings.TrimSpace(cfg.TraceID),
		IdempotencyKey:  idempotencyKey,
		Subject:         strings.TrimSpace(cfg.Subject),
		SpecVersion:     DefaultEventSpecVersion,
		DataContentType: ProtobufContentType,
		DataSchema:      strings.TrimSpace(cfg.DataSchema),
	}
}

func StableEventID(prefix string, parts ...string) string {
	return strings.TrimSpace(prefix) + hashx.StableParts(parts...)
}
