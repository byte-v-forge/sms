package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/eventoutbox"
)

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
