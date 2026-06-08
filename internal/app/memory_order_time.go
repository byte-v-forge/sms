package app

import (
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func memoryOrderAge(order core.Order, now time.Time) time.Duration {
	updatedAt := memoryOrderUpdatedAt(order)
	if updatedAt.IsZero() || now.Before(updatedAt) {
		return 0
	}
	return now.Sub(updatedAt)
}

func memoryOrderUpdatedAt(order core.Order) time.Time {
	for _, value := range []time.Time{order.UpdatedAt, order.ExpiresAt, order.AcquiredAt} {
		if !value.IsZero() {
			return value
		}
	}
	return time.Time{}
}
