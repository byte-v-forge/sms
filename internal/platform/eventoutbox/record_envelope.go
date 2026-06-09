package eventoutbox

import (
	"fmt"
	"strings"

	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
	"google.golang.org/protobuf/proto"
)

func recordFromEnvelope(envelope *commonv1.EventEnvelope) (Record, error) {
	metadata := envelope.GetMetadata()
	if metadata == nil || strings.TrimSpace(metadata.GetId()) == "" {
		return Record{}, ErrMissingEventID
	}
	payload, err := proto.Marshal(envelope)
	if err != nil {
		return Record{}, fmt.Errorf("marshal event outbox envelope: %w", err)
	}
	return Record{
		EventID:        metadata.GetId(),
		Subject:        envelope.GetSubject(),
		EventName:      metadata.GetType(),
		IdempotencyKey: metadata.GetIdempotencyKey(),
		Envelope:       payload,
	}, nil
}
