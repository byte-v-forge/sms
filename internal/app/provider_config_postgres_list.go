package app

import (
	"context"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
)

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
	return scanProviderConfigs(rows)
}
