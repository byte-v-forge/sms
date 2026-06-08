package fivesim

import (
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func defaultProviderPolicy() core.ProviderPolicy {
	return core.ProviderPolicy{
		OrderTTL:     20 * time.Minute,
		PollInterval: 5 * time.Second,
	}
}
