package eventbus

import (
	"fmt"
	"strings"

	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
	"google.golang.org/protobuf/proto"
)

func NewEnvelope(message Message) (*commonv1.EventEnvelope, error) {
	subject := strings.TrimSpace(message.Subject)
	if subject == "" {
		return nil, ErrEmptySubject
	}
	if message.Event == nil {
		return nil, ErrEmptyEvent
	}
	if err := ValidateMetadata(message.Metadata); err != nil {
		return nil, err
	}
	payload, err := proto.Marshal(message.Event)
	if err != nil {
		return nil, fmt.Errorf("marshal event payload: %w", err)
	}
	return &commonv1.EventEnvelope{
		Metadata:        message.Metadata,
		Subject:         subject,
		PayloadType:     string(message.Event.ProtoReflect().Descriptor().FullName()),
		Payload:         payload,
		DataContentType: ProtobufContentType,
		Extensions:      cloneExtensions(message.Extensions),
	}, nil
}

func cloneExtensions(values map[string]string) map[string]string {
	if len(values) == 0 {
		return nil
	}
	out := make(map[string]string, len(values))
	for key, value := range values {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		out[key] = strings.TrimSpace(value)
	}
	return out
}
