package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/eventoutbox"
)

type CompositeOrderStore struct {
	active  OrderStore
	history OrderStore
}

func NewCompositeOrderStore(active OrderStore, history OrderStore) *CompositeOrderStore {
	return &CompositeOrderStore{active: active, history: history}
}

func (s *CompositeOrderStore) Save(ctx context.Context, order core.Order, events ...eventoutbox.Record) error {
	if err := s.history.Save(ctx, order, events...); err != nil {
		return err
	}
	return s.active.Save(ctx, order)
}

func (s *CompositeOrderStore) Update(ctx context.Context, order core.Order, events ...eventoutbox.Record) error {
	if err := s.history.Update(ctx, order, events...); err != nil {
		return err
	}
	return s.active.Update(ctx, order)
}

func (s *CompositeOrderStore) RecordCode(ctx context.Context, order core.Order, code core.SMSCode, events ...eventoutbox.Record) error {
	if err := s.history.RecordCode(ctx, order, code, events...); err != nil {
		return err
	}
	return s.active.Update(ctx, order)
}

func (s *CompositeOrderStore) CodeSecretExists(ctx context.Context, orderID string, secretID string) (bool, error) {
	if s == nil || s.history == nil {
		return false, nil
	}
	return s.history.CodeSecretExists(ctx, orderID, secretID)
}

func (s *CompositeOrderStore) Get(ctx context.Context, orderID string) (core.Order, error) {
	order, err := s.active.Get(ctx, orderID)
	if err == nil {
		return order, nil
	}
	return s.history.Get(ctx, orderID)
}
