package eventbusadapter

import (
	"errors"

	"github.com/byte-v-forge/sms/internal/core"
)

func coreErrorRetryableUnless(err error, nonRetryableCodes ...core.ErrorCode) bool {
	var smsErr *core.Error
	if !errors.As(err, &smsErr) {
		return true
	}
	for _, code := range nonRetryableCodes {
		if smsErr.Code == code {
			return false
		}
	}
	return smsErr.Retryable
}
