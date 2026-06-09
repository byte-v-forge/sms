package app

import smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"

func PublicError(err error) *smsv1.SmsError {
	smsErr := publicCoreError(err)
	if smsErr == nil {
		return nil
	}
	return &smsv1.SmsError{Code: PublicErrorCode(smsErr.Code), Message: smsErr.Message, Retryable: smsErr.Retryable}
}
