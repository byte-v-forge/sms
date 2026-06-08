package eventoutbox

import (
	"context"
	"errors"
	"time"
)

const (
	StatusPending   = "PENDING"
	StatusPublished = "PUBLISHED"
	StatusDiscarded = "DISCARDED"
)

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
