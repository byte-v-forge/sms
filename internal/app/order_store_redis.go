package app

import (
	"context"
	"encoding/json"
	"time"

	"github.com/byte-v-forge/common-lib/eventoutbox"
	"github.com/byte-v-forge/common-lib/redisx"
	"github.com/byte-v-forge/sms/internal/core"
)

const terminalOrderTTL = 30 * time.Minute

type RedisOrderStore struct {
	store *redisx.StringStore
	clock core.Clock
}

func NewRedisOrderStore(store *redisx.StringStore, clock core.Clock) *RedisOrderStore {
	if clock == nil { clock = SystemClock{} }
	return &RedisOrderStore{store: store, clock: clock}
}

func (s *RedisOrderStore) Save(ctx context.Context, order core.Order, _ ...eventoutbox.Record) error { return s.save(ctx, order) }
func (s *RedisOrderStore) Update(ctx context.Context, order core.Order, _ ...eventoutbox.Record) error { return s.save(ctx, order) }

func (s *RedisOrderStore) Get(ctx context.Context, orderID string) (core.Order, error) {
	value, ok, err := s.store.Load(ctx, orderID)
	if err != nil { return core.Order{}, err }
	if !ok { return core.Order{}, core.NewError(core.CodeOrderNotFound, "order not found", false) }
	var order core.Order
	if err := json.Unmarshal([]byte(value), &order); err != nil { return core.Order{}, err }
	return order, nil
}

func (s *RedisOrderStore) save(ctx context.Context, order core.Order) error {
	payload, err := json.Marshal(order)
	if err != nil { return err }
	return s.store.SaveTTL(ctx, order.ID, string(payload), s.ttl(order))
}

func (s *RedisOrderStore) ttl(order core.Order) time.Duration {
	if order.Status.IsFinal() { return terminalOrderTTL }
	if !order.ExpiresAt.IsZero() {
		if ttl := order.ExpiresAt.Sub(s.clock.Now()); ttl > 0 { return ttl }
	}
	return s.store.DefaultTTL()
}

type CompositeOrderStore struct {
	active  OrderStore
	history OrderStore
}

func NewCompositeOrderStore(active OrderStore, history OrderStore) *CompositeOrderStore {
	return &CompositeOrderStore{active: active, history: history}
}

func (s *CompositeOrderStore) Save(ctx context.Context, order core.Order, events ...eventoutbox.Record) error {
	if err := s.history.Save(ctx, order, events...); err != nil { return err }
	return s.active.Save(ctx, order)
}

func (s *CompositeOrderStore) Update(ctx context.Context, order core.Order, events ...eventoutbox.Record) error {
	if err := s.history.Update(ctx, order, events...); err != nil { return err }
	return s.active.Update(ctx, order)
}

func (s *CompositeOrderStore) Get(ctx context.Context, orderID string) (core.Order, error) {
	order, err := s.active.Get(ctx, orderID)
	if err == nil { return order, nil }
	return s.history.Get(ctx, orderID)
}
