package app

import (
	"context"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
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
	s.applySaveMetadata(config, now)
	if err := validateProviderConfigCredential(config); err != nil {
		return nil, err
	}
	config.UpdatedAt = cloneTimestamp(now)
	markProviderConfigCredentialState(config)
	s.configs[config.GetProviderKey()] = cloneProviderConfig(config)
	return cloneProviderConfig(config), nil
}
