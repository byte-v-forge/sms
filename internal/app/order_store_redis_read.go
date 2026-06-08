package app

import (
	"context"
	"encoding/json"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *RedisOrderStore) Get(ctx context.Context, orderID string) (core.Order, error) {
	value, ok, err := s.store.Load(ctx, orderID)
	if err != nil {
		return core.Order{}, err
	}
	if !ok {
		return core.Order{}, core.NewError(core.CodeOrderNotFound, "order not found", false)
	}
	var order core.Order
	if err := json.Unmarshal([]byte(value), &order); err != nil {
		return core.Order{}, err
	}
	return order, nil
}
