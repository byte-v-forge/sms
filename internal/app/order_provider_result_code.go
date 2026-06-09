package app

import (
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func codeFromProviderResult(result core.ProviderCodeResult, fallbackReceivedAt time.Time) core.SMSCode {
	receivedAt := result.ReceivedAt
	if receivedAt.IsZero() {
		receivedAt = fallbackReceivedAt
	}
	return core.SMSCode{
		Value:       result.Code,
		MessageText: result.MessageText,
		ReceivedAt:  receivedAt,
	}
}
