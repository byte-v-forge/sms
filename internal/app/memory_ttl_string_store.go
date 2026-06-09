package app

import (
	"sync"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

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
