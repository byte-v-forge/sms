package eventbus

import (
	"fmt"
	"strings"
)

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
