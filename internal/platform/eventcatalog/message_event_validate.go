package eventcatalog

import (
	"fmt"
	"strings"

	"github.com/byte-v-forge/sms/internal/platform/eventbus"
	"google.golang.org/protobuf/proto"
)

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
