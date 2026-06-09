package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *MemoryOrderStore) Get(ctx context.Context, orderID string) (core.Order, error) {
	if err := ctx.Err(); err != nil {
		return core.Order{}, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	order, ok := s.orders[orderID]
	if !ok {
		return core.Order{}, core.NewError(core.CodeOrderNotFound, "order not found", false)
	}
	return cloneOrder(order), nil
}
