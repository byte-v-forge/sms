package app

import (
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func shouldQueueEarlyCancelRetry(err *core.Error, policy core.ProviderPolicy) bool {
	return err != nil && err.Code == core.CodeCancelNotAllowed && err.Retryable && policy.EarlyCancelRetryAfter > 0
}

func earlyCancelRetryAt(order core.Order, policy core.ProviderPolicy, now time.Time) time.Time {
	if !order.AcquiredAt.IsZero() && policy.EarlyCancelRetryAfter > 0 {
		retryAt := order.AcquiredAt.Add(policy.EarlyCancelRetryAfter)
		if retryAt.After(now) {
			return retryAt
		}
	}
	delay := policy.PollInterval
	if delay <= 0 {
		delay = 5 * time.Second
	}
	return now.Add(delay)
}
