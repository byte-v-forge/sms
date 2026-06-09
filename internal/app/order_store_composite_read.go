package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *CompositeOrderStore) Get(ctx context.Context, orderID string) (core.Order, error) {
	order, err := s.active.Get(ctx, orderID)
	if err == nil {
		return order, nil
	}
	return s.history.Get(ctx, orderID)
}
