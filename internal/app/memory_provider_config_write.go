package app

import (
	"context"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *MemoryProviderConfigStore) UpsertProviderConfig(ctx context.Context, input *smsinternalv1.SmsProviderConfig) (*smsinternalv1.SmsProviderConfig, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	config, err := s.prepareForSave(input)
	if err != nil {
		return nil, err
	}
	now := timestamppb.New(s.clock.Now())
	s.mu.Lock()
	defer s.mu.Unlock()
	existing := s.configs[config.GetProviderKey()]
	if existing != nil {
		config.CreatedAt = cloneTimestamp(existing.GetCreatedAt())
		if config.GetCredentialSecret() == "" {
			config.CredentialSecret = existing.GetCredentialSecret()
		}
	} else {
		config.CreatedAt = cloneTimestamp(now)
	}
	if err := validateProviderConfigCredential(config); err != nil {
		return nil, err
	}
	config.UpdatedAt = cloneTimestamp(now)
	markProviderConfigCredentialState(config)
	s.configs[config.GetProviderKey()] = cloneProviderConfig(config)
	return cloneProviderConfig(config), nil
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

func (s *MemoryProviderConfigStore) prepareForSave(input *smsinternalv1.SmsProviderConfig) (*smsinternalv1.SmsProviderConfig, error) {
	config, err := normalizeProviderConfigInput(input)
	if err != nil {
		return nil, err
	}
	if err := validateProviderConfigSupported(s.providers, config.GetProviderKey()); err != nil {
		return nil, err
	}
	return config, nil
}

func cloneTimestamp(ts *timestamppb.Timestamp) *timestamppb.Timestamp {
	if ts == nil {
		return nil
	}
	return timestamppb.New(ts.AsTime())
}
