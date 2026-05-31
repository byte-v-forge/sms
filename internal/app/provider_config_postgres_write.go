package app

import (
	"context"
	"strings"

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
	config := cloneProviderConfig(input)
	config.ProviderKey = normalizeProviderKey(config.GetProviderKey())
	if config.GetProviderKey() == "" {
		return nil, core.NewError(core.CodeValidationFailed, "provider_key is required", false)
	}
	if !supportedProviderKey(config.GetProviderKey()) {
		return nil, core.NewError(core.CodeUnsupportedOperation, "unsupported sms provider", false)
	}
	config.CredentialSecret = strings.TrimSpace(config.GetCredentialSecret())
	if config.GetCredentialSecret() == "" {
		existing, err := s.GetProviderConfig(ctx, config.GetProviderKey())
		if err == nil {
			config.CredentialSecret = existing.GetCredentialSecret()
		}
	}
	if config.GetEnabled() && config.GetCredentialSecret() == "" {
		return nil, core.NewError(core.CodeValidationFailed, "credential_secret is required for enabled sms provider", false)
	}
	config.CredentialSecretSet = strings.TrimSpace(config.GetCredentialSecret()) != ""
	return config, nil
}
