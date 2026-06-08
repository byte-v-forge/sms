package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *MemoryProviderConfigStore) DeleteProviderConfig(ctx context.Context, providerKey string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	providerKey = normalizeProviderKey(providerKey)
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.configs[providerKey] == nil {
		return core.NewError(core.CodeRouteNotFound, "sms provider config not found", false)
	}
	delete(s.configs, providerKey)
	return nil
}
