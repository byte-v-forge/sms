package fivesim

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func mapError(statusCode int, text string) error {
	normalized := strings.ToLower(strings.TrimSpace(text))
	switch {
	case statusCode == http.StatusUnauthorized:
		return core.NewError(core.CodeUpstreamRejected, "5sim credential rejected", false)
	case strings.Contains(normalized, "order not found"), strings.Contains(normalized, "record not found"):
		return core.NewError(core.CodeOrderNotFound, "5sim order not found", false)
	case strings.Contains(normalized, "no free phones"):
		return core.NewError(core.CodeNoNumberAvailable, "no 5sim phone available", true)
	case strings.Contains(normalized, "not enough user balance"), strings.Contains(normalized, "insufficient"):
		return core.NewError(core.CodeInsufficientBalance, "5sim balance is insufficient", false)
	case strings.Contains(normalized, "order expired"):
		return core.NewError(core.CodeOrderExpired, "5sim order expired", false)
	case strings.Contains(normalized, "order has sms"):
		return core.NewError(core.CodeCancelNotAllowed, "5sim order already has sms", false)
	case strings.Contains(normalized, "bad country"), strings.Contains(normalized, "bad operator"), strings.Contains(normalized, "bad product"),
		strings.Contains(normalized, "select country"), strings.Contains(normalized, "select operator"), strings.Contains(normalized, "select product"),
		strings.Contains(normalized, "product is empty"):
		return core.NewError(core.CodeValidationFailed, "5sim route parameters are invalid", false)
	case statusCode >= 500:
		return core.NewError(core.CodeSupplyUnavailable, "5sim supply is unavailable", true)
	case text == "":
		return core.NewError(core.CodeUpstreamRejected, fmt.Sprintf("5sim http status %d", statusCode), statusCode >= 500)
	default:
		return core.NewError(core.CodeUpstreamRejected, "5sim upstream rejected request", false)
	}
}
