package app

import (
	"context"
	"errors"
	"strings"
	"time"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PostgresProviderConfigStore struct{ pool *pgxpool.Pool }

func NewPostgresProviderConfigStore(ctx context.Context, dsn string) (*PostgresProviderConfigStore, error) {
	pool, err := newRequiredPostgresPool(ctx, dsn)
	if err != nil {
		return nil, err
	}
	store := &PostgresProviderConfigStore{pool: pool}
	if err := validatePostgresTables(ctx, pool, "sms_provider_configs"); err != nil {
		pool.Close()
		return nil, err
	}
	return store, nil
}

func (s *PostgresProviderConfigStore) Close() {
	if s != nil && s.pool != nil {
		s.pool.Close()
	}
}

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

func cloneProviderConfig(input *smsinternalv1.SmsProviderConfig) *smsinternalv1.SmsProviderConfig {
	if input == nil {
		return &smsinternalv1.SmsProviderConfig{}
	}
	return proto.Clone(input).(*smsinternalv1.SmsProviderConfig)
}

func defaultProviderCapabilities(providerKey string) *smsinternalv1.SmsProviderCapabilities {
	if plugin, ok := smsProviderPluginByKey(providerKey); ok {
		return plugin.Capabilities()
	}
	return &smsinternalv1.SmsProviderCapabilities{}
}

func supportedProviderKey(providerKey string) bool {
	_, ok := smsProviderPluginByKey(providerKey)
	return ok
}
