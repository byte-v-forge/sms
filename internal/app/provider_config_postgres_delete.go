package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

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
