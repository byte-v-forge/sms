package eventoutbox

import (
	"fmt"
	"strings"
)

func postgresIdentifier(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", ErrInvalidTableName
	}
	for _, part := range strings.Split(value, ".") {
		if !validPostgresIdentifierPart(part) {
			return "", fmt.Errorf("%w: %s", ErrInvalidTableName, value)
		}
	}
	return value, nil
}

func validPostgresIdentifierPart(value string) bool {
	if value == "" {
		return false
	}
	for i, r := range value {
		valid := r == '_' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || i > 0 && r >= '0' && r <= '9'
		if !valid {
			return false
		}
	}
	return true
}
