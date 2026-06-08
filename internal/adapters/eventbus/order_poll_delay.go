package eventbusadapter

import (
	"context"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func (w *OrderPollWorker) pollDelay(ctx context.Context, order core.Order) time.Duration {
	if order.ID == "" {
		return defaultPollRetryDelay
	}
	delay := w.service.PollInterval(ctx, order)
	if delay <= 0 {
		return defaultPollRetryDelay
	}
	return delay
}
