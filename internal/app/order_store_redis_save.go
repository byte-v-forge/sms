package app

import (
	"context"
	"encoding/json"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *RedisOrderStore) save(ctx context.Context, order core.Order) error {
	payload, err := json.Marshal(order)
	if err != nil {
		return err
	}
	return s.store.SaveTTL(ctx, order.ID, string(payload), s.ttl(order))
}
