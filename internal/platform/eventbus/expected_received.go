package eventbus

import (
	"fmt"
	"strings"
)

func (expected ExpectedMessage) ValidateReceived(received ReceivedMessage) error {
	if expected.IsZero() {
		return nil
	}
	if received.Envelope == nil {
		return ErrEmptyEnvelope
	}
	if err := expected.validateSubject(received); err != nil {
		return err
	}
	if err := expected.validateMetadata(received); err != nil {
		return err
	}
	if err := expected.validatePayloadType(received.Envelope.GetPayloadType()); err != nil {
		return err
	}
	return nil
}

func (expected ExpectedMessage) validateSubject(received ReceivedMessage) error {
	subject := strings.TrimSpace(expected.Subject)
	if subject == "" {
		return nil
	}
	receivedSubject := strings.TrimSpace(received.Subject)
	if receivedSubject != "" && receivedSubject != subject {
		return fmt.Errorf("%w: expected %s, got %s", ErrMismatchedSubject, subject, receivedSubject)
	}
	envelopeSubject := strings.TrimSpace(received.Envelope.GetSubject())
	if envelopeSubject == "" {
		return ErrEmptySubject
	}
	if envelopeSubject != subject {
		return fmt.Errorf("%w: expected %s, got %s", ErrMismatchedSubject, subject, envelopeSubject)
	}
	return nil
}

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
