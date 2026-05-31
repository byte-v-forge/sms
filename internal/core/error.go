package core

import "fmt"

type ErrorCode string

const (
	CodeValidationFailed      ErrorCode = "validation_failed"
	CodeRouteNotFound         ErrorCode = "route_not_found"
	CodeOrderNotFound         ErrorCode = "order_not_found"
	CodeOrderAlreadyFinalized ErrorCode = "order_already_finalized"
	CodeNoNumberAvailable     ErrorCode = "no_number_available"
	CodeRateLimited           ErrorCode = "rate_limited"
	CodeSupplyUnavailable     ErrorCode = "supply_unavailable"
	CodeUpstreamRejected      ErrorCode = "upstream_rejected"
	CodeTimeout               ErrorCode = "timeout"
	CodeUnsupportedOperation  ErrorCode = "unsupported_operation"
	CodeOrderExpired          ErrorCode = "order_expired"
	CodeCancelNotAllowed      ErrorCode = "cancel_not_allowed"
	CodeInsufficientBalance   ErrorCode = "insufficient_balance"
	CodeInternal              ErrorCode = "internal"
)

type Error struct {
	Code      ErrorCode
	Message   string
	Retryable bool
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Message == "" {
		return string(e.Code)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewError(code ErrorCode, message string, retryable bool) *Error {
	return &Error{Code: code, Message: message, Retryable: retryable}
}
