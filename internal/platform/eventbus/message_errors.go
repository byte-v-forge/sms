package eventbus

import "errors"

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
