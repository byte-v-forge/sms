package app

import (
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func orderRequestExpiresAt(now time.Time, policy core.ProviderPolicy, lease time.Duration) time.Time {
	policy = policy.WithDefaults()
	ttl := policy.OrderTTL
	if lease > 0 && lease < ttl {
		ttl = lease
	}
	if ttl <= 0 {
		return time.Time{}
	}
	return now.Add(ttl)
}

func remainingLease(now time.Time, expiresAt time.Time) time.Duration {
	if expiresAt.IsZero() {
		return 0
	}
	if !expiresAt.After(now) {
		return 0
	}
	return expiresAt.Sub(now)
}
