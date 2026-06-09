package app

import (
	"errors"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
	"github.com/jackc/pgx/v5"
)

func scanProviderConfig(row pgx.Row) (*smsinternalv1.SmsProviderConfig, error) {
	var record providerConfigRecord
	if err := row.Scan(&record.providerKey, &record.enabled, &record.credentialSecret, &record.createdAt, &record.updatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, core.NewError(core.CodeRouteNotFound, "sms provider config not found", false)
		}
		return nil, err
	}
	return record.toProto(), nil
}
