package fivesim

import (
	"time"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/stringx"
)

func smsHasContent(message sms) bool {
	return message.Code != "" || message.Text != ""
}

func smsCodeResult(message sms) core.ProviderCodeResult {
	return core.ProviderCodeResult{
		Status:      core.StatusCodeReceived,
		Code:        message.Code,
		MessageText: message.Text,
		ReceivedAt:  smsReceivedAt(message),
	}
}

func latestSMS(messages []sms) sms {
	if len(messages) == 0 {
		return sms{}
	}
	latest := messages[0]
	latestAt := smsReceivedAt(latest)
	for _, message := range messages[1:] {
		messageAt := smsReceivedAt(message)
		if latestAt.IsZero() || messageAt.After(latestAt) {
			latest = message
			latestAt = messageAt
		}
	}
	return latest
}

func smsReceivedAt(message sms) time.Time {
	return parseTime(stringx.FirstNonEmpty(message.Date, message.CreatedAt))
}
