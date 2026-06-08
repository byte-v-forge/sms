package handlerapi

import (
	"encoding/json"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/jsonx"
	"github.com/byte-v-forge/sms/internal/platform/stringx"
	"github.com/byte-v-forge/sms/internal/providers/providerhttp"
)

func MapTextError(text string) error {
	text = strings.TrimSpace(text)
	code, message := normalizeHandlerAPIError(text)
	switch {
	case text == "":
		return core.NewError(core.CodeUpstreamRejected, "empty upstream response", true)
	case code == "BAD_KEY":
		return core.NewError(core.CodeUpstreamRejected, "provider credential rejected", false)
	case code == "BAD_ACTION":
		return core.NewError(core.CodeUnsupportedOperation, "provider action rejected", false)
	case code == "BAD_SERVICE", code == "BAD_COUNTRY", code == "BAD_STATUS", code == "WRONG_EXCEPTION_PHONE", code == "WRONG_ACTIVATION_ID":
		return core.NewError(core.CodeValidationFailed, message, false)
	case code == "NO_ACTIVATION":
		return core.NewError(core.CodeOrderNotFound, "upstream order not found", false)
	case code == "NO_BALANCE":
		return core.NewError(core.CodeInsufficientBalance, "provider balance is insufficient", false)
	case code == "NO_NUMBERS", code == "NO_NUMBER", strings.Contains(text, "NO_NUMBERS"):
		return core.NewError(core.CodeNoNumberAvailable, "no upstream number available", true)
	case code == "EARLY_CANCEL_DENIED":
		return core.NewError(core.CodeCancelNotAllowed, "upstream denied early cancel", true)
	case code == "WRONG_MAX_PRICE", code == "ERROR_SQL", code == "ERROR_SQL25", code == "SERVER_ERROR":
		return core.NewError(core.CodeSupplyUnavailable, message, true)
	case code == "BANNED", code == "CHANNELS_LIMIT":
		return core.NewError(core.CodeSupplyUnavailable, message, false)
	case code == "SERVICE_NOT_AVAILABLE", code == "NOT_AVAILABLE", code == "WHATSAPP_NOT_AVAILABLE":
		return core.NewError(core.CodeSupplyUnavailable, message, true)
	default:
		return core.NewError(core.CodeUpstreamRejected, message, false)
	}
}

type handlerAPIErrorPayload struct {
	Title   string                     `json:"title"`
	Detail  string                     `json:"detail"`
	Details string                     `json:"details"`
	Info    map[string]json.RawMessage `json:"info"`
}

func normalizeHandlerAPIError(text string) (string, string) {
	text = strings.TrimSpace(text)
	if text == "" {
		return "", ""
	}
	if code, message, ok := parseHandlerAPIJSONError(text); ok {
		return code, message
	}
	code := text
	if idx := strings.Index(code, ":"); idx >= 0 {
		code = strings.TrimSpace(code[:idx])
	}
	return code, truncateHandlerAPIErrorMessage(text)
}

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

func truncateHandlerAPIErrorMessage(message string) string {
	return providerhttp.SanitizeErrorText(message)
}
