package app

import (
	"context"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
)

func (s *ProviderAdminService) GetProviderConfig(ctx context.Context, providerKey string) (*smsinternalv1.SmsProviderConfig, error) {
	config, err := s.configs.GetProviderConfig(ctx, providerKey)
	if err != nil {
		return nil, err
	}
	return RedactProviderConfig(config), nil
}

func (s *ProviderAdminService) ListProviderConfigs(ctx context.Context, includeDisabled bool, providerKey string) ([]*smsinternalv1.SmsProviderConfig, error) {
	configs, err := s.configs.ListProviderConfigs(ctx, includeDisabled, providerKey)
	if err != nil {
		return nil, err
	}
	for index, config := range configs {
		configs[index] = RedactProviderConfig(config)
	}
	return configs, nil
}
