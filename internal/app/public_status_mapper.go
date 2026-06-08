package app

import (
	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

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
