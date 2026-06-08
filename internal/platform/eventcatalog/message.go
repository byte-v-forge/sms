package eventcatalog

import (
	"errors"
	"fmt"
	"strings"

	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
	"github.com/byte-v-forge/sms/internal/platform/eventbus"
	"google.golang.org/protobuf/proto"
)

var (
	ErrEmptyDefinitionSubject      = errors.New("event catalog subject is required")
	ErrEmptyDefinitionEventName    = errors.New("event catalog event_name is required")
	ErrEmptyDefinitionEventVersion = errors.New("event catalog event_version is required")
	ErrEmptyDefinitionPayloadType  = errors.New("event catalog payload_type is required")
	ErrMismatchedEventName         = errors.New("event context event_name does not match catalog definition")
	ErrMismatchedEventVersion      = errors.New("event context event_version does not match catalog definition")
	ErrMismatchedPayloadType       = errors.New("event payload type does not match catalog definition")
)

func (definition Definition) NewMessage(
	event proto.Message,
	metadata *commonv1.EventMetadata,
	attributes map[string]string,
) (eventbus.Message, error) {
	if err := definition.ValidateEvent(event); err != nil {
		return eventbus.Message{}, err
	}
	if err := definition.ValidateMetadata(metadata); err != nil {
		return eventbus.Message{}, err
	}
	return eventbus.Message{
		Subject:    strings.TrimSpace(definition.Subject),
		Event:      event,
		Metadata:   metadata,
		Extensions: attributes,
	}, nil
}

func (definition Definition) ValidateEvent(event proto.Message) error {
	if strings.TrimSpace(definition.Subject) == "" {
		return ErrEmptyDefinitionSubject
	}
	if strings.TrimSpace(definition.PayloadType) == "" {
		return ErrEmptyDefinitionPayloadType
	}
	if event == nil {
		return eventbus.ErrEmptyEvent
	}
	actualType := string(event.ProtoReflect().Descriptor().FullName())
	if actualType != strings.TrimSpace(definition.PayloadType) {
		return fmt.Errorf("%w: expected %s, got %s", ErrMismatchedPayloadType, definition.PayloadType, actualType)
	}
	return nil
}

func (definition Definition) ValidateMetadata(metadata *commonv1.EventMetadata) error {
	if strings.TrimSpace(definition.EventName) == "" {
		return ErrEmptyDefinitionEventName
	}
	if strings.TrimSpace(definition.EventVersion) == "" {
		return ErrEmptyDefinitionEventVersion
	}
	if err := eventbus.ValidateMetadata(metadata); err != nil {
		return err
	}
	if strings.TrimSpace(metadata.GetType()) != strings.TrimSpace(definition.EventName) {
		return fmt.Errorf("%w: expected %s, got %s", ErrMismatchedEventName, definition.EventName, metadata.GetType())
	}
	if strings.TrimSpace(metadata.GetVersion()) != strings.TrimSpace(definition.EventVersion) {
		return fmt.Errorf("%w: expected %s, got %s", ErrMismatchedEventVersion, definition.EventVersion, metadata.GetVersion())
	}
	return nil
}
