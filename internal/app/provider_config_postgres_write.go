package app

import (
	"context"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
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
