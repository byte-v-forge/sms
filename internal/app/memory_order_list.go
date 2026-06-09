package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/pagex"
)

func (s *MemoryOrderStore) List(ctx context.Context, includeFinal bool, limit int) ([]core.Order, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	limit = pagex.NormalizeLimit(limit, pagex.DefaultLimit, pagex.MaxLimit)
	s.mu.RLock()
	defer s.mu.RUnlock()
	orders := s.filteredOrdersLocked(includeFinal)
	sortOrdersByUpdatedAt(orders)
	if len(orders) > limit {
		return orders[:limit], nil
	}
	return orders, nil
}
