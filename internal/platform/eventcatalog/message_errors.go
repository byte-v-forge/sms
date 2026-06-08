package eventcatalog

import "errors"

var (
	ErrEmptyDefinitionSubject      = errors.New("event catalog subject is required")
	ErrEmptyDefinitionEventName    = errors.New("event catalog event_name is required")
	ErrEmptyDefinitionEventVersion = errors.New("event catalog event_version is required")
	ErrEmptyDefinitionPayloadType  = errors.New("event catalog payload_type is required")
	ErrMismatchedEventName         = errors.New("event context event_name does not match catalog definition")
	ErrMismatchedEventVersion      = errors.New("event context event_version does not match catalog definition")
	ErrMismatchedPayloadType       = errors.New("event payload type does not match catalog definition")
)
