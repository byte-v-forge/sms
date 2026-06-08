package eventbus

import (
	"errors"
	"fmt"
	"strings"

	"google.golang.org/protobuf/proto"
)

var (
	ErrEmptyEnvelope          = errors.New("event envelope is required")
	ErrEmptyPayloadType       = errors.New("event payload type is required")
	ErrMismatchedSubject      = errors.New("event subject does not match expected event")
	ErrMismatchedEventType    = errors.New("event metadata type does not match expected event")
	ErrMismatchedEventVersion = errors.New("event metadata version does not match expected event")
	ErrMismatchedPayloadType  = errors.New("event payload type does not match expected event")
)

type ExpectedMessage struct {
	Subject      string
	EventName    string
	EventVersion string
	PayloadType  string
}

func (expected ExpectedMessage) IsZero() bool {
	return strings.TrimSpace(expected.Subject) == "" &&
		strings.TrimSpace(expected.EventName) == "" &&
		strings.TrimSpace(expected.EventVersion) == "" &&
		strings.TrimSpace(expected.PayloadType) == ""
}

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

func (expected ExpectedMessage) ValidateEvent(event proto.Message) error {
	payloadType := strings.TrimSpace(expected.PayloadType)
	if payloadType == "" {
		return nil
	}
	if event == nil {
		return ErrEmptyEvent
	}
	actualType := string(event.ProtoReflect().Descriptor().FullName())
	if actualType != payloadType {
		return fmt.Errorf("%w: expected %s, got %s", ErrMismatchedPayloadType, payloadType, actualType)
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

func (expected ExpectedMessage) validatePayloadType(actual string) error {
	payloadType := strings.TrimSpace(expected.PayloadType)
	if payloadType == "" {
		return nil
	}
	actual = strings.TrimSpace(actual)
	if actual == "" {
		return ErrEmptyPayloadType
	}
	if actual != payloadType {
		return fmt.Errorf("%w: expected %s, got %s", ErrMismatchedPayloadType, payloadType, actual)
	}
	return nil
}
