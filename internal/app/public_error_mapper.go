package app

import (
	"errors"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

func PublicError(err error) *smsv1.SmsError {
	if err == nil {
		return nil
	}
	var smsErr *core.Error
	if errors.As(err, &smsErr) {
		if smsErr == nil {
			return nil
		}
	} else {
		smsErr = runtimeCoreError(err)
	}
	if smsErr == nil {
		smsErr = core.NewError(core.CodeInternal, "sms service request failed", false)
	}
	if smsErr == nil {
		return nil
	}
	return &smsv1.SmsError{Code: PublicErrorCode(smsErr.Code), Message: smsErr.Message, Retryable: smsErr.Retryable}
}

func PublicErrorCode(code core.ErrorCode) smsv1.SmsErrorCode {
	switch code {
	case core.CodeValidationFailed:
		return smsv1.SmsErrorCode_SMS_ERROR_CODE_VALIDATION_FAILED
	case core.CodeRouteNotFound:
		return smsv1.SmsErrorCode_SMS_ERROR_CODE_ROUTE_NOT_FOUND
	case core.CodeOrderNotFound:
		return smsv1.SmsErrorCode_SMS_ERROR_CODE_ORDER_NOT_FOUND
	case core.CodeOrderAlreadyFinalized:
		return smsv1.SmsErrorCode_SMS_ERROR_CODE_ORDER_ALREADY_FINALIZED
	case core.CodeNoNumberAvailable:
		return smsv1.SmsErrorCode_SMS_ERROR_CODE_NO_NUMBER_AVAILABLE
	case core.CodeRateLimited:
		return smsv1.SmsErrorCode_SMS_ERROR_CODE_RATE_LIMITED
	case core.CodeSupplyUnavailable:
		return smsv1.SmsErrorCode_SMS_ERROR_CODE_SUPPLY_UNAVAILABLE
	case core.CodeUpstreamRejected:
		return smsv1.SmsErrorCode_SMS_ERROR_CODE_UPSTREAM_REJECTED
	case core.CodeTimeout:
		return smsv1.SmsErrorCode_SMS_ERROR_CODE_TIMEOUT
	case core.CodeUnsupportedOperation:
		return smsv1.SmsErrorCode_SMS_ERROR_CODE_UNSUPPORTED_OPERATION
	case core.CodeOrderExpired:
		return smsv1.SmsErrorCode_SMS_ERROR_CODE_ORDER_EXPIRED
	case core.CodeCancelNotAllowed:
		return smsv1.SmsErrorCode_SMS_ERROR_CODE_CANCEL_NOT_ALLOWED
	case core.CodeInsufficientBalance:
		return smsv1.SmsErrorCode_SMS_ERROR_CODE_INSUFFICIENT_BALANCE
	case core.CodeInternal:
		return smsv1.SmsErrorCode_SMS_ERROR_CODE_INTERNAL
	default:
		return smsv1.SmsErrorCode_SMS_ERROR_CODE_UNSPECIFIED
	}
}
