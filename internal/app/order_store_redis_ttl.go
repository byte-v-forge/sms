package app

import (
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

const terminalOrderTTL = 30 * time.Minute

func (s *RedisOrderStore) ttl(order core.Order) time.Duration {
	if order.Status.IsFinal() {
		return terminalOrderTTL
	}
	if !order.ExpiresAt.IsZero() {
		if ttl := order.ExpiresAt.Sub(s.clock.Now()); ttl > 0 {
			return ttl
		}
	}
	return s.store.DefaultTTL()
}
