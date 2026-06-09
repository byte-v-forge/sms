package app

import (
	"context"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

func (s *MemoryProviderConfigStore) GetProviderConfig(ctx context.Context, providerKey string) (*smsinternalv1.SmsProviderConfig, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	providerKey = normalizeProviderKey(providerKey)
	if providerKey == "" {
		return nil, core.NewError(core.CodeValidationFailed, "provider_key is required", false)
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	config := s.configs[providerKey]
	if config == nil {
		return nil, core.NewError(core.CodeRouteNotFound, "sms provider config not found", false)
	}
	return cloneProviderConfig(config), nil
}
