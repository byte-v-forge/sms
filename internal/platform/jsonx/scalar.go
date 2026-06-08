package jsonx

import (
	"encoding/json"
	"strconv"
	"strings"
)

// Scalar converts a JSON string, number, or null scalar into a trimmed string.
func Scalar(raw json.RawMessage) string {
	if len(raw) == 0 || string(raw) == "null" {
		return ""
	}
	var text string
	if err := json.Unmarshal(raw, &text); err == nil {
		return strings.TrimSpace(text)
	}
	var number json.Number
	if err := json.Unmarshal(raw, &number); err == nil {
		return strings.TrimSpace(number.String())
	}
	var floatValue float64
	if err := json.Unmarshal(raw, &floatValue); err == nil {
		return strconv.FormatFloat(floatValue, 'f', -1, 64)
	}
	return strings.TrimSpace(strings.Trim(string(raw), "\""))
}

// FirstScalar returns the first non-empty scalar from the provided JSON values.
func FirstScalar(values ...json.RawMessage) string {
	for _, value := range values {
		if scalar := Scalar(value); scalar != "" {
			return scalar
		}
	}
	return ""
}

// Int parses a JSON scalar as an integer and returns zero when the value is empty or malformed.
func Int(raw json.RawMessage) int {
	value := Scalar(raw)
	if value == "" {
		return 0
	}
	parsed, _ := strconv.Atoi(value)
	return parsed
}
