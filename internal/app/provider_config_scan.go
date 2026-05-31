package app

import (
	"errors"
	"strings"
	"time"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
	"github.com/jackc/pgx/v5"
	"google.golang.org/protobuf/types/known/timestamppb"
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

type providerConfigRecord struct {
	providerKey      string
	enabled          bool
	credentialSecret string
	createdAt        time.Time
	updatedAt        time.Time
}

func (r providerConfigRecord) toProto() *smsinternalv1.SmsProviderConfig {
	return &smsinternalv1.SmsProviderConfig{
		ProviderKey:         normalizeProviderKey(r.providerKey),
		Enabled:             r.enabled,
		CredentialSecret:    r.credentialSecret,
		CredentialSecretSet: strings.TrimSpace(r.credentialSecret) != "",
		CreatedAt:           timestamppb.New(r.createdAt),
		UpdatedAt:           timestamppb.New(r.updatedAt),
	}
}

func providerConfigColumns() string {
	return `provider_key, enabled, credential_secret, created_at, updated_at`
}
