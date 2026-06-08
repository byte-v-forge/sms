package eventoutbox

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
	"github.com/byte-v-forge/sms/internal/platform/eventbus"
	"github.com/byte-v-forge/sms/internal/platform/eventcatalog"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

const (
	StatusPending   = "PENDING"
	StatusPublished = "PUBLISHED"
	StatusDiscarded = "DISCARDED"
)

const defaultPublishTimeout = 10 * time.Second

var (
	ErrMissingEventID = errors.New("event outbox event_id is required")
	ErrNilPublisher   = errors.New("event outbox publisher is required")
	ErrNilUpdates     = errors.New("event outbox updates is required")
)

type Record struct {
	EventID        string
	Subject        string
	EventName      string
	IdempotencyKey string
	Envelope       []byte
}

type Row struct {
	EventID      string
	Envelope     []byte
	AttemptCount int32
}

type Updates interface {
	MarkPublished(ctx context.Context, eventID string, publishedAt int64) error
	MarkRetry(ctx context.Context, eventID string, attemptCount int32, nextAttemptAt int64, lastError string, updatedAt int64) error
	MarkDiscarded(ctx context.Context, eventID string, lastError string, updatedAt int64) error
}

type PublishOptions struct {
	PublishTimeout time.Duration
	RetryDelay     func(int32) time.Duration
	Now            func() time.Time
}

func NewRecord(message eventbus.Message) (Record, error) {
	envelope, err := eventbus.NewEnvelope(message)
	if err != nil {
		return Record{}, err
	}
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

func NewRecordFor(
	definition eventcatalog.Definition,
	event proto.Message,
	metadata *commonv1.EventMetadata,
	attributes map[string]string,
) (Record, error) {
	message, err := definition.NewMessage(event, metadata, attributes)
	if err != nil {
		return Record{}, err
	}
	return NewRecord(message)
}

func PublishRows(ctx context.Context, publisher eventbus.Publisher, rows []Row, updates Updates, options PublishOptions) (int, error) {
	if publisher == nil {
		return 0, ErrNilPublisher
	}
	if updates == nil {
		return 0, ErrNilUpdates
	}
	published := 0
	for _, row := range rows {
		if ctx.Err() != nil {
			return published, ctx.Err()
		}
		message, err := MessageFromEnvelope(row.Envelope)
		now := optionNow(options).Unix()
		if err != nil {
			if updateErr := updates.MarkDiscarded(ctx, row.EventID, TruncateError(err), now); updateErr != nil {
				return published, updateErr
			}
			continue
		}
		publishCtx, cancel := context.WithTimeout(ctx, publishTimeout(options))
		_, err = publisher.Publish(publishCtx, message)
		cancel()
		if err != nil {
			nextAttempt := row.AttemptCount + 1
			nextAttemptAt := optionNow(options).Add(retryDelay(options, nextAttempt)).Unix()
			if updateErr := updates.MarkRetry(ctx, row.EventID, nextAttempt, nextAttemptAt, TruncateError(err), optionNow(options).Unix()); updateErr != nil {
				return published, updateErr
			}
			continue
		}
		if updateErr := updates.MarkPublished(ctx, row.EventID, optionNow(options).Unix()); updateErr != nil {
			return published, updateErr
		}
		published++
	}
	return published, nil
}

func MessageFromEnvelope(payload []byte) (eventbus.Message, error) {
	envelope := &commonv1.EventEnvelope{}
	if err := proto.Unmarshal(payload, envelope); err != nil {
		return eventbus.Message{}, fmt.Errorf("decode event outbox envelope: %w", err)
	}
	messageType, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(envelope.GetPayloadType()))
	if err != nil {
		return eventbus.Message{}, fmt.Errorf("resolve event outbox payload type %s: %w", envelope.GetPayloadType(), err)
	}
	message := messageType.New().Interface()
	if err := proto.Unmarshal(envelope.GetPayload(), message); err != nil {
		return eventbus.Message{}, fmt.Errorf("decode event outbox payload %s: %w", envelope.GetPayloadType(), err)
	}
	return eventbus.Message{
		Subject:    envelope.GetSubject(),
		Event:      message,
		Metadata:   envelope.GetMetadata(),
		Extensions: envelope.GetExtensions(),
	}, nil
}

func DefaultRetryDelay(attempt int32) time.Duration {
	switch {
	case attempt <= 1:
		return 5 * time.Second
	case attempt == 2:
		return 15 * time.Second
	case attempt == 3:
		return 30 * time.Second
	case attempt <= 6:
		return time.Minute
	default:
		return 5 * time.Minute
	}
}

func TruncateError(err error) string {
	if err == nil {
		return ""
	}
	message := strings.TrimSpace(err.Error())
	if len(message) <= 1000 {
		return message
	}
	return message[:1000]
}

func publishTimeout(options PublishOptions) time.Duration {
	if options.PublishTimeout > 0 {
		return options.PublishTimeout
	}
	return defaultPublishTimeout
}

func retryDelay(options PublishOptions, attempt int32) time.Duration {
	if options.RetryDelay != nil {
		return options.RetryDelay(attempt)
	}
	return DefaultRetryDelay(attempt)
}

func optionNow(options PublishOptions) time.Time {
	if options.Now != nil {
		return options.Now()
	}
	return time.Now()
}
