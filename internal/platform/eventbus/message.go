package eventbus

import (
	"context"
	"errors"
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
