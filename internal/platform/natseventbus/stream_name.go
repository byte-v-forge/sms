package natseventbus

import "strings"

func normalizedStreamName(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return DefaultStream
	}
	return value
}
