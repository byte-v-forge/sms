package app

import (
	"context"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

func (s *PostgresProviderConfigStore) GetProviderConfig(ctx context.Context, providerKey string) (*smsinternalv1.SmsProviderConfig, error) {
	providerKey = normalizeProviderKey(providerKey)
	if providerKey == "" {
		return nil, core.NewError(core.CodeValidationFailed, "provider_key is required", false)
	}
	row := s.pool.QueryRow(ctx, `SELECT `+providerConfigColumns()+` FROM sms_provider_configs WHERE provider_key = $1`, providerKey)
	return scanProviderConfig(row)
}

func (s *PostgresProviderConfigStore) ListProviderConfigs(ctx context.Context, includeDisabled bool, providerKey string) ([]*smsinternalv1.SmsProviderConfig, error) {
	providerKey = normalizeProviderKey(providerKey)
	rows, err := s.pool.Query(ctx, `
SELECT `+providerConfigColumns()+`
FROM sms_provider_configs
WHERE ($1 OR enabled) AND ($2 = '' OR provider_key = $2)
ORDER BY provider_key ASC
`, includeDisabled, providerKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	configs := []*smsinternalv1.SmsProviderConfig{}
	for rows.Next() {
		config, err := scanProviderConfig(rows)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}
	return configs, rows.Err()
}

func (s *PostgresProviderConfigStore) GetEnabledProviderConfig(ctx context.Context, providerKey string, _ core.Target) (*smsinternalv1.SmsProviderConfig, error) {
	configs, err := s.ListProviderConfigs(ctx, false, providerKey)
	if err != nil {
		return nil, err
	}
	if len(configs) == 0 {
		return nil, core.NewError(core.CodeRouteNotFound, "no enabled sms provider config", false)
	}
	return configs[0], nil
}
