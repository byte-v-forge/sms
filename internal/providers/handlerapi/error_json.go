package handlerapi

import (
	"encoding/json"
	"strings"

	"github.com/byte-v-forge/sms/internal/platform/jsonx"
	"github.com/byte-v-forge/sms/internal/platform/stringx"
)

func parseHandlerAPIJSONError(text string) (string, string, bool) {
	var payload handlerAPIErrorPayload
	if err := json.Unmarshal([]byte(text), &payload); err != nil {
		return "", "", false
	}
	code := strings.TrimSpace(payload.Title)
	if code == "" {
		return "", "", false
	}
	parts := []string{code}
	if details := stringx.FirstNonEmpty(payload.Details, payload.Detail); details != "" {
		parts = append(parts, details)
	}
	for _, key := range []string{"min", "max"} {
		if value := jsonx.Scalar(payload.Info[key]); value != "" {
			parts = append(parts, key+"="+value)
		}
	}
	return code, truncateHandlerAPIErrorMessage(strings.Join(parts, ": ")), true
}
