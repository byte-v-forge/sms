package fivesim

import (
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/stringx"
)

func orderToCodeResult(payload order) core.ProviderCodeResult {
	latestSMS := latestSMS(payload.SMS)
	if latestSMS.Code != "" || latestSMS.Text != "" {
		receivedAt := parseTime(stringx.FirstNonEmpty(latestSMS.Date, latestSMS.CreatedAt))
		return core.ProviderCodeResult{
			Status:      core.StatusCodeReceived,
			Code:        latestSMS.Code,
			MessageText: latestSMS.Text,
			ReceivedAt:  receivedAt,
		}
	}
	switch strings.ToUpper(strings.TrimSpace(payload.Status)) {
	case "CANCELED":
		return core.ProviderCodeResult{Status: core.StatusCanceled}
	case "TIMEOUT":
		return core.ProviderCodeResult{Status: core.StatusExpired}
	case "FINISHED":
		return core.ProviderCodeResult{Status: core.StatusCompleted}
	case "BANNED":
		return core.ProviderCodeResult{Status: core.StatusFailed}
	default:
		return core.ProviderCodeResult{Status: core.StatusPendingCode}
	}
}

func latestSMS(messages []sms) sms {
	if len(messages) == 0 {
		return sms{}
	}
	latest := messages[0]
	latestAt := parseTime(stringx.FirstNonEmpty(latest.Date, latest.CreatedAt))
	for _, message := range messages[1:] {
		messageAt := parseTime(stringx.FirstNonEmpty(message.Date, message.CreatedAt))
		if latestAt.IsZero() || messageAt.After(latestAt) {
			latest = message
			latestAt = messageAt
		}
	}
	return latest
}
