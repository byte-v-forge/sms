package jsonx

import (
	"encoding/json"
	"strings"
)

// Scalar converts a JSON string, number, or null scalar into a trimmed string.
func Scalar(raw json.RawMessage) string {
	if len(raw) == 0 || string(raw) == "null" {
		return ""
	}
	if text, ok := scalarString(raw); ok {
		return text
	}
	if number, ok := scalarNumber(raw); ok {
		return number
	}
	return strings.TrimSpace(strings.Trim(string(raw), "\""))
}
