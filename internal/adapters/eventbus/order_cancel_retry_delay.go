package eventbusadapter

import (
	"errors"
	"time"

	"github.com/byte-v-forge/sms/internal/app"
	"github.com/byte-v-forge/sms/internal/core"
)

func cancelRetryDelay(err error, now time.Time) (time.Duration, bool) {
	var retryErr *app.CancelRetryError
	if !errors.As(err, &retryErr) {
		return 0, false
	}
	delay := retryErr.RetryAt.Sub(now)
	if retryErr.RetryAt.IsZero() || delay <= 0 {
		delay = defaultCancelRetryDelay
	}
	return delay, true
}

func (w *OrderCancelWorker) cancelDelay(order core.Order) time.Duration {
	if order.ID != "" && !order.CancelAllowedAt.IsZero() {
		delay := time.Until(order.CancelAllowedAt)
		if delay > 0 {
			return delay
		}
	}
	return defaultCancelRetryDelay
}
