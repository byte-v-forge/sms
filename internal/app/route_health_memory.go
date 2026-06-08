package app

import (
	"sync"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

type MemoryRouteHealthStore struct {
	mu       sync.Mutex
	failures map[string]memoryRouteFailure
	disabled map[string]time.Time
	clock    core.Clock
}

type memoryRouteFailure struct {
	count     int
	expiresAt time.Time
}

func NewMemoryRouteHealthStore(clock core.Clock) *MemoryRouteHealthStore {
	if clock == nil {
		clock = SystemClock{}
	}
	return &MemoryRouteHealthStore{
		failures: map[string]memoryRouteFailure{},
		disabled: map[string]time.Time{},
		clock:    clock,
	}
}

func routeHealthExpired(now time.Time, expiresAt time.Time) bool {
	return !expiresAt.IsZero() && !now.Before(expiresAt)
}
