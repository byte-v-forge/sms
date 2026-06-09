package fivesim

import "github.com/byte-v-forge/sms/internal/core"

func orderToCodeResult(payload order) core.ProviderCodeResult {
	latestSMS := latestSMS(payload.SMS)
	if smsHasContent(latestSMS) {
		return smsCodeResult(latestSMS)
	}
	return orderStatusResult(payload.Status)
}
