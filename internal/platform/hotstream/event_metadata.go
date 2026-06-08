package hotstream

import (
	"time"

	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func eventMetadata(cfg EventConfig) *commonv1.EventMetadata {
	occurredAt := cfg.OccurredAt
	if occurredAt.IsZero() {
		occurredAt = time.Now()
	}
	return &commonv1.EventMetadata{
		Id:              cleanEventText(cfg.EventID),
		Type:            cleanEventText(cfg.EventType),
		Version:         "v1",
		Time:            timestamppb.New(occurredAt),
		Source:          cleanEventText(cfg.SourceService),
		CorrelationId:   cleanEventText(cfg.CorrelationID),
		TraceId:         cleanEventText(cfg.TraceID),
		SpecVersion:     "1.0",
		DataContentType: DataContentType,
	}
}
