package app

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

type TTLStringStore interface {
	DefaultTTL() time.Duration
	Load(ctx context.Context, key string) (string, bool, error)
	SaveTTL(ctx context.Context, key string, value string, ttl time.Duration) error
}

type MemoryTTLStringStore struct {
	mu       sync.RWMutex
	values   map[string]memoryTTLValue
	ttl      time.Duration
	clock    core.Clock
	keyspace string
}

type memoryTTLValue struct {
	value     string
	expiresAt time.Time
}

func NewMemoryTTLStringStore(prefix string, ttl time.Duration, clock core.Clock) *MemoryTTLStringStore {
	if ttl <= 0 {
		ttl = 30 * time.Minute
	}
	if clock == nil {
		clock = SystemClock{}
	}
	return &MemoryTTLStringStore{values: map[string]memoryTTLValue{}, ttl: ttl, clock: clock, keyspace: strings.Trim(strings.TrimSpace(prefix), ":")}
}

func (s *MemoryTTLStringStore) DefaultTTL() time.Duration {
	if s == nil {
		return 0
	}
	return s.ttl
}

func (s *MemoryTTLStringStore) Load(ctx context.Context, key string) (string, bool, error) {
	if err := ctx.Err(); err != nil {
		return "", false, err
	}
	normalized := s.key(key)
	if normalized == "" {
		return "", false, nil
	}
	s.mu.RLock()
	item, ok := s.values[normalized]
	s.mu.RUnlock()
	if !ok {
		return "", false, nil
	}
	if !item.expiresAt.IsZero() && !s.clock.Now().Before(item.expiresAt) {
		s.mu.Lock()
		delete(s.values, normalized)
		s.mu.Unlock()
		return "", false, nil
	}
	return item.value, true, nil
}

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

func (s *MemoryTTLStringStore) key(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if s.keyspace == "" {
		return value
	}
	return s.keyspace + ":" + value
}
