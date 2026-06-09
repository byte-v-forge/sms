package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/eventoutbox"
)

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
