package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/eventoutbox"
)

func (s *OrderService) saveOrder(ctx context.Context, order core.Order, records ...eventoutbox.Record) error {
	if err := s.store.Save(ctx, order, records...); err != nil {
		return err
	}
	s.publishOrder(ctx, order)
	return nil
}

func (s *OrderService) updateOrder(ctx context.Context, order core.Order, records ...eventoutbox.Record) error {
	if err := s.store.Update(ctx, order, records...); err != nil {
		return err
	}
	s.publishOrder(ctx, order)
	return nil
}

func (s *OrderService) recordCode(ctx context.Context, order core.Order, code core.SMSCode, records ...eventoutbox.Record) error {
	if err := s.store.RecordCode(ctx, order, code, records...); err != nil {
		return err
	}
	s.publishOrder(ctx, order)
	return nil
}
