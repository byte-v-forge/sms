package eventcatalog

import (
	"fmt"
	"strings"

	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
	"github.com/byte-v-forge/sms/internal/platform/eventbus"
)

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
