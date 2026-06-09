package jsonx

import (
	"encoding/json"
	"strconv"
)

// Int parses a JSON scalar as an integer and returns zero when the value is empty or malformed.
func Int(raw json.RawMessage) int {
	value := Scalar(raw)
	if value == "" {
		return 0
	}
	parsed, _ := strconv.Atoi(value)
	return parsed
}
