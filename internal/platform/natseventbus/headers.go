package natseventbus

import (
	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
	"github.com/nats-io/nats.go"
)

func envelopeHeaders(envelope *commonv1.EventEnvelope) nats.Header {
	headers := nats.Header{}
	if envelope == nil {
		return headers
	}
	headers.Set("Bvf-Event-Subject", envelope.GetSubject())
	headers.Set("Bvf-Event-Type", envelope.GetPayloadType())
	headers.Set("Content-Type", envelope.GetDataContentType())
	if metadata := envelope.GetMetadata(); metadata != nil {
		headers.Set("Bvf-Event-Id", metadata.GetId())
		headers.Set("Bvf-Event-Name", metadata.GetType())
		headers.Set("Bvf-Event-Version", metadata.GetVersion())
		headers.Set("Bvf-Correlation-Id", metadata.GetCorrelationId())
		headers.Set("Bvf-Trace-Id", metadata.GetTraceId())
		headers.Set("Bvf-Idempotency-Key", metadata.GetIdempotencyKey())
	}
	return headers
}
