package app

import (
	"context"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *MemoryTTLStringStore) SaveTTL(ctx context.Context, key string, value string, ttl time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	normalized := s.key(key)
	if normalized == "" {
		return core.NewError(core.CodeValidationFailed, "memory string store key is required", false)
	}
	if ttl <= 0 {
		ttl = s.ttl
	}
	expiresAt := time.Time{}
	if ttl > 0 {
		expiresAt = s.clock.Now().Add(ttl)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[normalized] = memoryTTLValue{value: value, expiresAt: expiresAt}
	return nil
}
