package app

import (
	"errors"
	"time"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/secretref"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func PublicOrder(order core.Order) *smsv1.SmsOrder {
	return &smsv1.SmsOrder{
		OrderId:   order.ID,
		RequestId: order.RequestID,
		Target: &smsv1.SmsTarget{
			ApplicationKey:     order.Target.ApplicationKey,
			CountryIso2:        order.Target.CountryISO2,
			CountryCallingCode: order.Target.CountryCallingCode,
		},
		PhoneNumber: &smsv1.PhoneNumber{
			E164Number:         order.PhoneNumber.E164,
			NationalNumber:     order.PhoneNumber.NationalNumber,
			CountryIso2:        order.PhoneNumber.CountryISO2,
			CountryCallingCode: order.PhoneNumber.CountryCallingCode,
		},
		Status:                   PublicOrderStatus(order.Status),
		Price:                    PublicMoney(order.Price),
		AcquiredAt:               PublicTime(order.AcquiredAt),
		ExpiresAt:                PublicTime(order.ExpiresAt),
		UpdatedAt:                PublicTime(order.UpdatedAt),
		LastError:                PublicError(order.LastError),
		CanRequestAdditionalCode: order.CanRequestAdditionalCode,
		CancelAllowedAt:          PublicTime(order.CancelAllowedAt),
	}
}

func PublicCode(code *core.SMSCode) *smsv1.SmsCode {
	if code == nil {
		return nil
	}
	return &smsv1.SmsCode{
		SecretRef:  secretref.Clone(code.SecretRef, "sms", "sms_otp"),
		ReceivedAt: PublicTime(code.ReceivedAt),
	}
}

func PublicMoney(money core.Money) *smsv1.DecimalMoney {
	if money.CurrencyCode == "" && money.AmountDecimal == "" {
		return nil
	}
	return &smsv1.DecimalMoney{CurrencyCode: money.CurrencyCode, AmountDecimal: money.AmountDecimal}
}

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

func PublicOrderStatus(status core.OrderStatus) smsv1.SmsOrderStatus {
	switch status {
	case core.StatusAcquireRequested:
		return smsv1.SmsOrderStatus_SMS_ORDER_STATUS_ACQUIRE_REQUESTED
	case core.StatusPendingCode:
		return smsv1.SmsOrderStatus_SMS_ORDER_STATUS_PENDING_CODE
	case core.StatusMessageSent:
		return smsv1.SmsOrderStatus_SMS_ORDER_STATUS_MESSAGE_SENT
	case core.StatusCodeReceived:
		return smsv1.SmsOrderStatus_SMS_ORDER_STATUS_CODE_RECEIVED
	case core.StatusAdditionalCodeRequested:
		return smsv1.SmsOrderStatus_SMS_ORDER_STATUS_ADDITIONAL_CODE_REQUESTED
	case core.StatusCompleted:
		return smsv1.SmsOrderStatus_SMS_ORDER_STATUS_COMPLETED
	case core.StatusCanceled:
		return smsv1.SmsOrderStatus_SMS_ORDER_STATUS_CANCELED
	case core.StatusExpired:
		return smsv1.SmsOrderStatus_SMS_ORDER_STATUS_EXPIRED
	case core.StatusFailed:
		return smsv1.SmsOrderStatus_SMS_ORDER_STATUS_FAILED
	default:
		return smsv1.SmsOrderStatus_SMS_ORDER_STATUS_UNSPECIFIED
	}
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

func PublicTime(value time.Time) *timestamppb.Timestamp {
	if value.IsZero() {
		return nil
	}
	return timestamppb.New(value)
}
