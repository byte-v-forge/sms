package eventbus

import (
	"fmt"
	"strings"
)

func (expected ExpectedMessage) validateMetadata(received ReceivedMessage) error {
	eventName := strings.TrimSpace(expected.EventName)
	eventVersion := strings.TrimSpace(expected.EventVersion)
	if eventName == "" && eventVersion == "" {
		return nil
	}
	metadata := received.Envelope.GetMetadata()
	if err := ValidateMetadata(metadata); err != nil {
		return err
	}
	if eventName != "" && strings.TrimSpace(metadata.GetType()) != eventName {
		return fmt.Errorf("%w: expected %s, got %s", ErrMismatchedEventType, eventName, metadata.GetType())
	}
	if eventVersion != "" && strings.TrimSpace(metadata.GetVersion()) != eventVersion {
		return fmt.Errorf("%w: expected %s, got %s", ErrMismatchedEventVersion, eventVersion, metadata.GetVersion())
	}
	return nil
}
