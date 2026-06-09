package app

import (
	"context"
	"sort"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
)

func (s *MemoryProviderConfigStore) ListProviderConfigs(ctx context.Context, includeDisabled bool, providerKey string) ([]*smsinternalv1.SmsProviderConfig, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	providerKey = normalizeProviderKey(providerKey)
	s.mu.RLock()
	defer s.mu.RUnlock()
	configs := s.matchingProviderConfigsLocked(includeDisabled, providerKey)
	sort.Slice(configs, func(i, j int) bool { return configs[i].GetProviderKey() < configs[j].GetProviderKey() })
	return configs, nil
}

func (s *MemoryProviderConfigStore) matchingProviderConfigsLocked(includeDisabled bool, providerKey string) []*smsinternalv1.SmsProviderConfig {
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
	return configs
}
