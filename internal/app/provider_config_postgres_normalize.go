package app

import (
	"context"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
)

func (s *PostgresProviderConfigStore) normalizeForSave(ctx context.Context, input *smsinternalv1.SmsProviderConfig) (*smsinternalv1.SmsProviderConfig, error) {
	config, err := normalizeProviderConfigInput(input)
	if err != nil {
		return nil, err
	}
	if err := validateProviderConfigSupported(s.providers, config.GetProviderKey()); err != nil {
		return nil, err
	}
	if err := s.fillExistingProviderCredential(ctx, config); err != nil {
		return nil, err
	}
	if err := validateProviderConfigCredential(config); err != nil {
		return nil, err
	}
	markProviderConfigCredentialState(config)
	return config, nil
}

func (s *PostgresProviderConfigStore) fillExistingProviderCredential(ctx context.Context, config *smsinternalv1.SmsProviderConfig) error {
	if config.GetCredentialSecret() != "" {
		return nil
	}
	existing, err := s.GetProviderConfig(ctx, config.GetProviderKey())
	if err == nil {
		config.CredentialSecret = existing.GetCredentialSecret()
	}
	return nil
}
