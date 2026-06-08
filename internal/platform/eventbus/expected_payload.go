package eventbus

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/proto"
)

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
