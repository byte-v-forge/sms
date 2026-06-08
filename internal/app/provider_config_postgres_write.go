package app

import (
	"context"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

func (s *PostgresProviderConfigStore) UpsertProviderConfig(ctx context.Context, input *smsinternalv1.SmsProviderConfig) (*smsinternalv1.SmsProviderConfig, error) {
	config, err := s.normalizeForSave(ctx, input)
	if err != nil {
		return nil, err
	}
	row := s.pool.QueryRow(ctx, `
INSERT INTO sms_provider_configs (provider_key, enabled, credential_secret)
VALUES ($1,$2,$3)
ON CONFLICT (provider_key) DO UPDATE SET
  enabled = EXCLUDED.enabled,
  credential_secret = EXCLUDED.credential_secret,
  updated_at = now()
RETURNING `+providerConfigColumns(), config.GetProviderKey(), config.GetEnabled(), config.GetCredentialSecret())
	return scanProviderConfig(row)
}

func (s *PostgresProviderConfigStore) DeleteProviderConfig(ctx context.Context, providerKey string) error {
	result, err := s.pool.Exec(ctx, `DELETE FROM sms_provider_configs WHERE provider_key = $1`, normalizeProviderKey(providerKey))
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return core.NewError(core.CodeRouteNotFound, "sms provider config not found", false)
	}
	return nil
}

func (s *PostgresProviderConfigStore) normalizeForSave(ctx context.Context, input *smsinternalv1.SmsProviderConfig) (*smsinternalv1.SmsProviderConfig, error) {
	config, err := normalizeProviderConfigInput(input)
	if err != nil {
		return nil, err
	}
	if err := validateProviderConfigSupported(s.providers, config.GetProviderKey()); err != nil {
		return nil, err
	}
	if config.GetCredentialSecret() == "" {
		existing, err := s.GetProviderConfig(ctx, config.GetProviderKey())
		if err == nil {
			config.CredentialSecret = existing.GetCredentialSecret()
		}
	}
	if err := validateProviderConfigCredential(config); err != nil {
		return nil, err
	}
	markProviderConfigCredentialState(config)
	return config, nil
}
