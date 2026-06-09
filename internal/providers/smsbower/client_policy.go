package smsbower

import (
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func defaultProviderPolicy() core.ProviderPolicy {
	return core.ProviderPolicy{
		OrderTTL:              25 * time.Minute,
		PollInterval:          5 * time.Second,
		EarlyCancelRetryAfter: 2 * time.Minute,
	}
}
