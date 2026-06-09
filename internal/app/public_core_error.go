package app

import (
	"errors"

	"github.com/byte-v-forge/sms/internal/core"
)

func publicCoreError(err error) *core.Error {
	if err == nil {
		return nil
	}
	var smsErr *core.Error
	if errors.As(err, &smsErr) {
		return smsErr
	}
	if smsErr = runtimeCoreError(err); smsErr != nil {
		return smsErr
	}
	return core.NewError(core.CodeInternal, "sms service request failed", false)
}
