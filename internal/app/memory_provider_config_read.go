package app

import (
	"context"
	"sort"

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

func (s *MemoryProviderConfigStore) ListProviderConfigs(ctx context.Context, includeDisabled bool, providerKey string) ([]*smsinternalv1.SmsProviderConfig, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	providerKey = normalizeProviderKey(providerKey)
	s.mu.RLock()
	defer s.mu.RUnlock()
	configs := make([]*smsinternalv1.SmsProviderConfig, 0, len(s.configs))
	for key, config := range s.configs {
		if providerKey != "" && key != providerKey {
			continue
		}
		if !includeDisabled && !config.GetEnabled() {
			continue
		}
		configs = append(configs, cloneProviderConfig(config))
	}
	sort.Slice(configs, func(i, j int) bool { return configs[i].GetProviderKey() < configs[j].GetProviderKey() })
	return configs, nil
}

func (s *MemoryProviderConfigStore) GetEnabledProviderConfig(ctx context.Context, providerKey string, _ core.Target) (*smsinternalv1.SmsProviderConfig, error) {
	configs, err := s.ListProviderConfigs(ctx, false, providerKey)
	if err != nil {
		return nil, err
	}
	if len(configs) == 0 {
		return nil, core.NewError(core.CodeRouteNotFound, "no enabled sms provider config", false)
	}
	return configs[0], nil
}
