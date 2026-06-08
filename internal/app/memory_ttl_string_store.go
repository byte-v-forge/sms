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
