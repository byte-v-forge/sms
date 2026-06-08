package app

import (
	"context"
	"sort"
	"strings"
	"sync"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
	providerspi "github.com/byte-v-forge/sms/internal/providers/spi"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (s *MemoryProviderConfigStore) UpsertProviderConfig(ctx context.Context, input *smsinternalv1.SmsProviderConfig) (*smsinternalv1.SmsProviderConfig, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	config, err := s.normalizeForSave(input)
	if err != nil {
		return nil, err
	}
	now := timestamppb.New(s.clock.Now())
	s.mu.Lock()
	defer s.mu.Unlock()
	if existing := s.configs[config.GetProviderKey()]; existing != nil {
		config.CreatedAt = cloneTimestamp(existing.GetCreatedAt())
		if strings.TrimSpace(config.GetCredentialSecret()) == "" {
			config.CredentialSecret = existing.GetCredentialSecret()
		}
	} else {
		config.CreatedAt = cloneTimestamp(now)
	}
	config.UpdatedAt = cloneTimestamp(now)
	config.CredentialSecretSet = strings.TrimSpace(config.GetCredentialSecret()) != ""
	s.configs[config.GetProviderKey()] = cloneProviderConfig(config)
	return cloneProviderConfig(config), nil
}

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

func (s *MemoryProviderConfigStore) normalizeForSave(input *smsinternalv1.SmsProviderConfig) (*smsinternalv1.SmsProviderConfig, error) {
	config := cloneProviderConfig(input)
	config.ProviderKey = normalizeProviderKey(config.GetProviderKey())
	if config.GetProviderKey() == "" {
		return nil, core.NewError(core.CodeValidationFailed, "provider_key is required", false)
	}
	if s.providers != nil && !s.providers.Supports(config.GetProviderKey()) {
		return nil, core.NewError(core.CodeUnsupportedOperation, "unsupported sms provider", false)
	}
	config.CredentialSecret = strings.TrimSpace(config.GetCredentialSecret())
	if config.GetCredentialSecret() == "" {
		s.mu.RLock()
		existing := s.configs[config.GetProviderKey()]
		if existing != nil {
			config.CredentialSecret = existing.GetCredentialSecret()
		}
		s.mu.RUnlock()
	}
	if config.GetEnabled() && config.GetCredentialSecret() == "" {
		return nil, core.NewError(core.CodeValidationFailed, "credential_secret is required for enabled sms provider", false)
	}
	return config, nil
}

func cloneTimestamp(ts *timestamppb.Timestamp) *timestamppb.Timestamp {
	if ts == nil {
		return nil
	}
	return timestamppb.New(ts.AsTime())
}
