package app

import "github.com/byte-v-forge/sms/internal/core"

func asCoreError(err error) *core.Error {
	if err == nil {
		return nil
	}
	if smsErr, ok := err.(*core.Error); ok {
		return smsErr
	}
	if smsErr := runtimeCoreError(err); smsErr != nil {
		return smsErr
	}
	return core.NewError(core.CodeInternal, "sms service request failed", false)
}
