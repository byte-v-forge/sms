package app

import (
	"sync"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
	providerspi "github.com/byte-v-forge/sms/internal/providers/spi"
)

type MemoryProviderConfigStore struct {
	mu        sync.RWMutex
	providers *providerspi.Registry
	configs   map[string]*smsinternalv1.SmsProviderConfig
	clock     core.Clock
}

func NewMemoryProviderConfigStore(providers *providerspi.Registry, clock core.Clock) *MemoryProviderConfigStore {
	if clock == nil {
		clock = SystemClock{}
	}
	return &MemoryProviderConfigStore{providers: providers, configs: map[string]*smsinternalv1.SmsProviderConfig{}, clock: clock}
}
