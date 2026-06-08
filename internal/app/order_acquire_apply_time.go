package app

import (
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

const providerClockSkewTolerance = 5 * time.Minute

func normalizedProviderAcquiredAt(providerAcquiredAt time.Time, now time.Time) time.Time {
	if providerAcquiredAt.IsZero() || providerAcquiredAt.After(now.Add(providerClockSkewTolerance)) {
		return now
	}
	return providerAcquiredAt
}

func normalizedProviderExpiresAt(orderExpiresAt time.Time, providerExpiresAt time.Time, acquiredAt time.Time, policy core.ProviderPolicy) time.Time {
	expiresAt := providerExpiresAt
	if expiresAt.IsZero() || expiresAt.Before(acquiredAt) {
		expiresAt = acquiredAt.Add(policy.OrderTTL)
	}
	if !orderExpiresAt.IsZero() && orderExpiresAt.Before(expiresAt) {
		expiresAt = orderExpiresAt
	}
	return expiresAt
}

func providerCancelAllowedAt(acquiredAt time.Time, policy core.ProviderPolicy) time.Time {
	if policy.CancelAllowedAfter <= 0 {
		return time.Time{}
	}
	return acquiredAt.Add(policy.CancelAllowedAfter)
}
