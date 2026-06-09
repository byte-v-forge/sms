package app

import (
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func prepareSyncedProviderOrder(order core.Order, now time.Time) core.Order {
	order.UpdatedAt = now
	order.LastError = nil
	order.CancelAllowedAt = time.Time{}
	return order
}
