package eventbus

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
	"google.golang.org/protobuf/proto"
)

const ProtobufContentType = "application/x-protobuf"

var (
	ErrEmptySubject        = errors.New("event subject is required")
	ErrEmptyEvent          = errors.New("event message is required")
	ErrEmptyPayload        = errors.New("event payload is required")
	ErrEmptyEventMetadata  = errors.New("event metadata is required")
	ErrEmptyEventID        = errors.New("event metadata id is required")
	ErrEmptyEventType      = errors.New("event metadata type is required")
	ErrEmptyEventVersion   = errors.New("event metadata version is required")
	ErrEmptySource         = errors.New("event metadata source is required")
	ErrEmptyIdempotencyKey = errors.New("event metadata idempotency_key is required")
	ErrEmptyEventTime      = errors.New("event metadata time is required")
)

type Message struct {
	Subject    string
	Event      proto.Message
	Metadata   *commonv1.EventMetadata
	Extensions map[string]string
}

type ReceivedMessage struct {
	Subject    string
	Envelope   *commonv1.EventEnvelope
	Extensions map[string]string
	Attempt    int32
	Ack        func(context.Context) error
	Nak        func(context.Context) error
	NakDelay   func(context.Context, time.Duration) error
	Term       func(context.Context) error
	DeadLetter func(context.Context, string) error
}

type PublishAck struct {
	Stream    string
	Sequence  uint64
	Duplicate bool
}

type Publisher interface {
	Publish(context.Context, Message) (PublishAck, error)
}

type Consumer interface {
	Fetch(context.Context, int) ([]ReceivedMessage, error)
}

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

func ValidateMetadata(metadata *commonv1.EventMetadata) error {
	if metadata == nil {
		return ErrEmptyEventMetadata
	}
	if strings.TrimSpace(metadata.GetId()) == "" {
		return ErrEmptyEventID
	}
	if strings.TrimSpace(metadata.GetType()) == "" {
		return ErrEmptyEventType
	}
	if strings.TrimSpace(metadata.GetVersion()) == "" {
		return ErrEmptyEventVersion
	}
	if metadata.GetTime() == nil || !metadata.GetTime().IsValid() {
		return ErrEmptyEventTime
	}
	if strings.TrimSpace(metadata.GetSource()) == "" {
		return ErrEmptySource
	}
	if strings.TrimSpace(metadata.GetIdempotencyKey()) == "" {
		return ErrEmptyIdempotencyKey
	}
	return nil
}

func UnmarshalPayload(message ReceivedMessage, event proto.Message) error {
	if event == nil {
		return ErrEmptyEvent
	}
	if message.Envelope == nil || len(message.Envelope.GetPayload()) == 0 {
		return ErrEmptyPayload
	}
	if err := proto.Unmarshal(message.Envelope.GetPayload(), event); err != nil {
		return fmt.Errorf("unmarshal event payload: %w", err)
	}
	return nil
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
