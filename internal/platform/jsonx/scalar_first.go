package jsonx

import "encoding/json"

// FirstScalar returns the first non-empty scalar from the provided JSON values.
func FirstScalar(values ...json.RawMessage) string {
	for _, value := range values {
		if scalar := Scalar(value); scalar != "" {
			return scalar
		}
	}
	return ""
}
