package app

import (
	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

var publicErrorCodes = map[core.ErrorCode]smsv1.SmsErrorCode{
	core.CodeValidationFailed:      smsv1.SmsErrorCode_SMS_ERROR_CODE_VALIDATION_FAILED,
	core.CodeRouteNotFound:         smsv1.SmsErrorCode_SMS_ERROR_CODE_ROUTE_NOT_FOUND,
	core.CodeOrderNotFound:         smsv1.SmsErrorCode_SMS_ERROR_CODE_ORDER_NOT_FOUND,
	core.CodeOrderAlreadyFinalized: smsv1.SmsErrorCode_SMS_ERROR_CODE_ORDER_ALREADY_FINALIZED,
	core.CodeNoNumberAvailable:     smsv1.SmsErrorCode_SMS_ERROR_CODE_NO_NUMBER_AVAILABLE,
	core.CodeRateLimited:           smsv1.SmsErrorCode_SMS_ERROR_CODE_RATE_LIMITED,
	core.CodeSupplyUnavailable:     smsv1.SmsErrorCode_SMS_ERROR_CODE_SUPPLY_UNAVAILABLE,
	core.CodeUpstreamRejected:      smsv1.SmsErrorCode_SMS_ERROR_CODE_UPSTREAM_REJECTED,
	core.CodeTimeout:               smsv1.SmsErrorCode_SMS_ERROR_CODE_TIMEOUT,
	core.CodeUnsupportedOperation:  smsv1.SmsErrorCode_SMS_ERROR_CODE_UNSUPPORTED_OPERATION,
	core.CodeOrderExpired:          smsv1.SmsErrorCode_SMS_ERROR_CODE_ORDER_EXPIRED,
	core.CodeCancelNotAllowed:      smsv1.SmsErrorCode_SMS_ERROR_CODE_CANCEL_NOT_ALLOWED,
	core.CodeInsufficientBalance:   smsv1.SmsErrorCode_SMS_ERROR_CODE_INSUFFICIENT_BALANCE,
	core.CodeInternal:              smsv1.SmsErrorCode_SMS_ERROR_CODE_INTERNAL,
}

func PublicErrorCode(code core.ErrorCode) smsv1.SmsErrorCode {
	if publicCode, ok := publicErrorCodes[code]; ok {
		return publicCode
	}
	return smsv1.SmsErrorCode_SMS_ERROR_CODE_UNSPECIFIED
}
