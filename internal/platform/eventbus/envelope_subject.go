package eventbus

import "strings"

func envelopeSubject(value string) (string, error) {
	subject := strings.TrimSpace(value)
	if subject == "" {
		return "", ErrEmptySubject
	}
	return subject, nil
}
