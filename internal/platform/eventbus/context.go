package eventbus

import (
	"strings"
	"time"

	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
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
	return &commonv1.EventMetadata{
		Id:              strings.TrimSpace(cfg.EventID),
		Type:            strings.TrimSpace(cfg.EventName),
		Version:         eventMetadataVersion(cfg.EventVersion),
		Time:            timestamppb.New(eventMetadataTime(cfg.OccurredAt)),
		Source:          strings.TrimSpace(cfg.SourceService),
		CorrelationId:   strings.TrimSpace(cfg.CorrelationID),
		TraceId:         strings.TrimSpace(cfg.TraceID),
		IdempotencyKey:  eventMetadataIdempotencyKey(cfg),
		Subject:         strings.TrimSpace(cfg.Subject),
		SpecVersion:     DefaultEventSpecVersion,
		DataContentType: ProtobufContentType,
		DataSchema:      strings.TrimSpace(cfg.DataSchema),
	}
}
