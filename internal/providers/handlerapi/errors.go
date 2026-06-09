package handlerapi

import (
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
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
	case strings.Contains(code, "OFFER_NOT_FOUND"):
		return core.NewError(core.CodeNoNumberAvailable, "upstream offer is no longer available", true)
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
