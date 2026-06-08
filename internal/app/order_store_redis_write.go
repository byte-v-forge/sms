package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/eventoutbox"
)

func (s *RedisOrderStore) Save(ctx context.Context, order core.Order, _ ...eventoutbox.Record) error {
	return s.save(ctx, order)
}

func (s *RedisOrderStore) Update(ctx context.Context, order core.Order, _ ...eventoutbox.Record) error {
	return s.save(ctx, order)
}

func (s *RedisOrderStore) RecordCode(ctx context.Context, order core.Order, _ core.SMSCode, _ ...eventoutbox.Record) error {
	return s.save(ctx, order)
}

func (s *RedisOrderStore) CodeSecretExists(context.Context, string, string) (bool, error) {
	return false, nil
}
