package app

import (
	"context"
	"sort"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/pagex"
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

func (s *MemoryOrderStore) List(ctx context.Context, includeFinal bool, limit int) ([]core.Order, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	limit = pagex.NormalizeLimit(limit, pagex.DefaultLimit, pagex.MaxLimit)
	s.mu.RLock()
	defer s.mu.RUnlock()
	orders := make([]core.Order, 0, len(s.orders))
	for _, order := range s.orders {
		if !includeFinal && order.Status.IsFinal() {
			continue
		}
		orders = append(orders, cloneOrder(order))
	}
	sort.Slice(orders, func(i, j int) bool {
		if orders[i].UpdatedAt.Equal(orders[j].UpdatedAt) {
			return orders[i].ID < orders[j].ID
		}
		return orders[i].UpdatedAt.After(orders[j].UpdatedAt)
	})
	if len(orders) > limit {
		return orders[:limit], nil
	}
	return orders, nil
}
