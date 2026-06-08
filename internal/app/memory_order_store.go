package app

import (
	"context"
	"sync"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/eventoutbox"
)

type MemoryOrderStore struct {
	mu     sync.RWMutex
	orders map[string]core.Order
	codes  map[string][]core.OrderCode
	clock  core.Clock
}

func NewMemoryOrderStore() *MemoryOrderStore {
	return &MemoryOrderStore{orders: map[string]core.Order{}, codes: map[string][]core.OrderCode{}, clock: SystemClock{}}
}

func (s *MemoryOrderStore) Save(ctx context.Context, order core.Order, _ ...eventoutbox.Record) error {
	return s.save(ctx, order)
}

func (s *MemoryOrderStore) Update(ctx context.Context, order core.Order, _ ...eventoutbox.Record) error {
	return s.save(ctx, order)
}

func (s *MemoryOrderStore) save(ctx context.Context, order core.Order) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[order.ID] = cloneOrder(order)
	s.pruneLocked(s.clock.Now())
	return nil
}
