package app

import (
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func NewMemoryTTLStringStore(prefix string, ttl time.Duration, clock core.Clock) *MemoryTTLStringStore {
	if ttl <= 0 {
		ttl = 30 * time.Minute
	}
	if clock == nil {
		clock = SystemClock{}
	}
	return &MemoryTTLStringStore{values: map[string]memoryTTLValue{}, ttl: ttl, clock: clock, keyspace: strings.Trim(strings.TrimSpace(prefix), ":")}
}
