package app

import (
	"context"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
)

func (s *ProviderAdminService) UpsertProviderConfig(ctx context.Context, config *smsinternalv1.SmsProviderConfig) (*smsinternalv1.SmsProviderConfig, error) {
	saved, err := s.configs.UpsertProviderConfig(ctx, config)
	if err != nil {
		return nil, err
	}
	s.publishProviderConfig(ctx, SMSProviderConfigUpdated, saved)
	return RedactProviderConfig(saved), nil
}

func (s *ProviderAdminService) DeleteProviderConfig(ctx context.Context, providerKey string) error {
	if err := s.configs.DeleteProviderConfig(ctx, providerKey); err != nil {
		return err
	}
	s.publishResource(ctx, SMSProviderConfigDeleted, SMSProviderConfigResource, providerKey, map[string]string{"provider_key": providerKey})
	return nil
}
