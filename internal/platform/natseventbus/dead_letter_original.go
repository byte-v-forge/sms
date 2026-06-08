package natseventbus

import commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"

type originalEventInfo struct {
	id            string
	eventName     string
	eventVersion  string
	source        string
	correlationID string
	traceID       string
}

func originalEvent(envelope *commonv1.EventEnvelope) originalEventInfo {
	metadata := envelope.GetMetadata()
	if metadata == nil {
		return originalEventInfo{}
	}
	return originalEventInfo{
		id:            metadata.GetId(),
		eventName:     metadata.GetType(),
		eventVersion:  metadata.GetVersion(),
		source:        metadata.GetSource(),
		correlationID: metadata.GetCorrelationId(),
		traceID:       metadata.GetTraceId(),
	}
}
