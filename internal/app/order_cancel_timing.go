package app

import (
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func normalizeProviderCancelTimes(order core.Order, policy core.ProviderPolicy, now time.Time) core.Order {
	if order.AcquiredAt.After(now.Add(providerClockSkewTolerance)) {
		order.AcquiredAt = now
		if policy.CancelAllowedAfter > 0 {
			order.AcquiredAt = now.Add(-policy.CancelAllowedAfter)
		}
	}
	if !order.CancelAllowedAt.IsZero() && !order.ExpiresAt.IsZero() && order.CancelAllowedAt.After(order.ExpiresAt) {
		order.CancelAllowedAt = time.Time{}
	}
	if !order.CancelAllowedAt.IsZero() && order.CancelAllowedAt.After(now.Add(providerClockSkewTolerance)) {
		order.CancelAllowedAt = order.AcquiredAt.Add(policy.CancelAllowedAfter)
	}
	return order
}
