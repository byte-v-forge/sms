package app

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PostgresProviderConfigStore struct {
	pool *pgxpool.Pool
}

func NewPostgresProviderConfigStore(ctx context.Context, dsn string) (*PostgresProviderConfigStore, error) {
	if strings.TrimSpace(dsn) == "" {
		return nil, errors.New("SMS_PG_DSN or PG_DSN is required")
	}
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	store := &PostgresProviderConfigStore{pool: pool}
	if err := store.EnsureSchema(ctx); err != nil {
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

func (s *PostgresProviderConfigStore) EnsureSchema(ctx context.Context) error {
	if err := s.ensureProviderConfigSchema(ctx); err != nil {
		return err
	}
	return s.ensureRouteProfileSchema(ctx)
}

func (s *PostgresProviderConfigStore) ensureProviderConfigSchema(ctx context.Context) error {
	_, err := s.pool.Exec(ctx, `
CREATE TABLE IF NOT EXISTS sms_provider_configs (
  provider_config_id text PRIMARY KEY,
  provider_key text NOT NULL,
  display_name text NOT NULL,
  enabled boolean NOT NULL DEFAULT true,
  api_endpoint text NOT NULL DEFAULT '',
  credential_secret text NOT NULL DEFAULT '',
  credential_secret_ref text NOT NULL DEFAULT '',
  proxy_ref text NOT NULL DEFAULT '',
  http_proxy text NOT NULL DEFAULT '',
  upstream_service_key text NOT NULL DEFAULT '',
  provider_country_id text NOT NULL DEFAULT '',
  default_application_key text NOT NULL DEFAULT '',
  default_country_iso2 text NOT NULL DEFAULT '',
  default_country_calling_code text NOT NULL DEFAULT '',
  default_min_price_currency text NOT NULL DEFAULT '',
  default_min_price_amount text NOT NULL DEFAULT '',
  default_max_price_currency text NOT NULL DEFAULT '',
  default_max_price_amount text NOT NULL DEFAULT '',
  policy_activation_ttl_seconds bigint NOT NULL DEFAULT 0,
  policy_poll_interval_seconds bigint NOT NULL DEFAULT 0,
  policy_cancel_allowed_after_seconds bigint NOT NULL DEFAULT 0,
  policy_early_cancel_retry_after_seconds bigint NOT NULL DEFAULT 0,
  policy_cancel_allowed_until_seconds bigint NOT NULL DEFAULT 0,
  capabilities jsonb NOT NULL DEFAULT '{}',
  labels jsonb NOT NULL DEFAULT '{}',
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_sms_provider_configs_enabled
  ON sms_provider_configs(provider_key, enabled, updated_at DESC);
ALTER TABLE sms_provider_configs ADD COLUMN IF NOT EXISTS policy_activation_ttl_seconds bigint NOT NULL DEFAULT 0;
ALTER TABLE sms_provider_configs ADD COLUMN IF NOT EXISTS policy_poll_interval_seconds bigint NOT NULL DEFAULT 0;
ALTER TABLE sms_provider_configs ADD COLUMN IF NOT EXISTS policy_cancel_allowed_after_seconds bigint NOT NULL DEFAULT 0;
ALTER TABLE sms_provider_configs ADD COLUMN IF NOT EXISTS policy_early_cancel_retry_after_seconds bigint NOT NULL DEFAULT 0;
ALTER TABLE sms_provider_configs ADD COLUMN IF NOT EXISTS policy_cancel_allowed_until_seconds bigint NOT NULL DEFAULT 0;
ALTER TABLE sms_provider_configs ADD COLUMN IF NOT EXISTS default_min_price_currency text NOT NULL DEFAULT '';
ALTER TABLE sms_provider_configs ADD COLUMN IF NOT EXISTS default_min_price_amount text NOT NULL DEFAULT '';
`)
	return err
}

func (s *PostgresProviderConfigStore) UpsertProviderConfig(ctx context.Context, input *smsinternalv1.SmsProviderConfig) (*smsinternalv1.SmsProviderConfig, error) {
	config, err := s.normalizeForSave(ctx, input)
	if err != nil {
		return nil, err
	}
	capabilities, err := json.Marshal(config.GetCapabilities())
	if err != nil {
		return nil, err
	}
	labels, err := json.Marshal(config.GetLabels())
	if err != nil {
		return nil, err
	}
	target := config.GetDefaultTarget()
	minPrice := target.GetMinPrice()
	maxPrice := target.GetMaxPrice()
	policy := providerPolicyFromConfig(config, defaultProviderPolicy(config.GetProviderKey())).WithDefaults()
	row := s.pool.QueryRow(ctx, `
INSERT INTO sms_provider_configs (
  provider_config_id, provider_key, display_name, enabled, api_endpoint,
  credential_secret, credential_secret_ref, proxy_ref, http_proxy,
  upstream_service_key, provider_country_id,
  default_application_key, default_country_iso2, default_country_calling_code,
  default_min_price_currency, default_min_price_amount,
  default_max_price_currency, default_max_price_amount,
  policy_activation_ttl_seconds, policy_poll_interval_seconds,
  policy_cancel_allowed_after_seconds, policy_early_cancel_retry_after_seconds,
  policy_cancel_allowed_until_seconds, capabilities, labels
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25)
ON CONFLICT (provider_config_id) DO UPDATE SET
  provider_key = EXCLUDED.provider_key,
  display_name = EXCLUDED.display_name,
  enabled = EXCLUDED.enabled,
  api_endpoint = EXCLUDED.api_endpoint,
  credential_secret = EXCLUDED.credential_secret,
  credential_secret_ref = EXCLUDED.credential_secret_ref,
  proxy_ref = EXCLUDED.proxy_ref,
  http_proxy = EXCLUDED.http_proxy,
  upstream_service_key = EXCLUDED.upstream_service_key,
  provider_country_id = EXCLUDED.provider_country_id,
  default_application_key = EXCLUDED.default_application_key,
  default_country_iso2 = EXCLUDED.default_country_iso2,
  default_country_calling_code = EXCLUDED.default_country_calling_code,
  default_min_price_currency = EXCLUDED.default_min_price_currency,
  default_min_price_amount = EXCLUDED.default_min_price_amount,
  default_max_price_currency = EXCLUDED.default_max_price_currency,
  default_max_price_amount = EXCLUDED.default_max_price_amount,
  policy_activation_ttl_seconds = EXCLUDED.policy_activation_ttl_seconds,
  policy_poll_interval_seconds = EXCLUDED.policy_poll_interval_seconds,
  policy_cancel_allowed_after_seconds = EXCLUDED.policy_cancel_allowed_after_seconds,
  policy_early_cancel_retry_after_seconds = EXCLUDED.policy_early_cancel_retry_after_seconds,
  policy_cancel_allowed_until_seconds = EXCLUDED.policy_cancel_allowed_until_seconds,
  capabilities = EXCLUDED.capabilities,
  labels = EXCLUDED.labels,
  updated_at = now()
`+selectProviderConfigSQL(),
		config.GetProviderConfigId(), config.GetProviderKey(), config.GetDisplayName(), config.GetEnabled(), config.GetApiEndpoint(),
		config.GetCredentialSecret(), config.GetCredentialSecretRef(), config.GetProxyRef(), config.GetHttpProxy(),
		config.GetUpstreamServiceKey(), config.GetProviderCountryId(),
		target.GetApplicationKey(), target.GetCountryIso2(), target.GetCountryCallingCode(),
		minPrice.GetCurrencyCode(), minPrice.GetAmountDecimal(),
		maxPrice.GetCurrencyCode(), maxPrice.GetAmountDecimal(),
		durationSeconds(policy.ActivationTTL), durationSeconds(policy.PollInterval),
		durationSeconds(policy.CancelAllowedAfter), durationSeconds(policy.EarlyCancelRetryAfter),
		durationSeconds(policy.CancelAllowedUntil), capabilities, labels,
	)
	return scanProviderConfig(row)
}

func (s *PostgresProviderConfigStore) GetProviderConfig(ctx context.Context, id string) (*smsinternalv1.SmsProviderConfig, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, core.NewError(core.CodeValidationFailed, "provider_config_id is required", false)
	}
	row := s.pool.QueryRow(ctx, `SELECT `+providerConfigColumns()+` FROM sms_provider_configs WHERE provider_config_id = $1`, id)
	return scanProviderConfig(row)
}

func (s *PostgresProviderConfigStore) ListProviderConfigs(ctx context.Context, includeDisabled bool, providerKey string) ([]*smsinternalv1.SmsProviderConfig, error) {
	providerKey = normalizeProviderKey(providerKey)
	rows, err := s.pool.Query(ctx, `
SELECT `+providerConfigColumns()+`
FROM sms_provider_configs
WHERE ($1 OR enabled) AND ($2 = '' OR provider_key = $2)
ORDER BY provider_key ASC, updated_at DESC, provider_config_id ASC
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

func (s *PostgresProviderConfigStore) DeleteProviderConfig(ctx context.Context, id string) error {
	result, err := s.pool.Exec(ctx, `DELETE FROM sms_provider_configs WHERE provider_config_id = $1`, strings.TrimSpace(id))
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return core.NewError(core.CodeRouteNotFound, "sms provider config not found", false)
	}
	return nil
}

func (s *PostgresProviderConfigStore) GetEnabledProviderConfig(ctx context.Context, providerKey string, target core.Target) (*smsinternalv1.SmsProviderConfig, error) {
	configs, err := s.ListProviderConfigs(ctx, false, providerKey)
	if err != nil {
		return nil, err
	}
	bestScore := -1
	var best *smsinternalv1.SmsProviderConfig
	for _, config := range configs {
		score := targetMatchScore(config.GetDefaultTarget(), target)
		if score < 0 {
			continue
		}
		if score > bestScore {
			bestScore = score
			best = config
		}
	}
	if best == nil {
		return nil, core.NewError(core.CodeRouteNotFound, "no enabled sms provider config for target", false)
	}
	return best, nil
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
	if strings.TrimSpace(config.GetProviderConfigId()) == "" {
		config.ProviderConfigId = config.GetProviderKey()
	}
	config.ProviderConfigId = strings.TrimSpace(config.GetProviderConfigId())
	config.DisplayName = strings.TrimSpace(config.GetDisplayName())
	if config.GetDisplayName() == "" {
		config.DisplayName = config.GetProviderKey()
	}
	config.ApiEndpoint = strings.TrimSpace(config.GetApiEndpoint())
	config.CredentialSecret = strings.TrimSpace(config.GetCredentialSecret())
	config.CredentialSecretRef = strings.TrimSpace(config.GetCredentialSecretRef())
	config.ProxyRef = strings.TrimSpace(config.GetProxyRef())
	config.HttpProxy = strings.TrimSpace(config.GetHttpProxy())
	config.UpstreamServiceKey = strings.TrimSpace(config.GetUpstreamServiceKey())
	config.ProviderCountryId = strings.TrimSpace(config.GetProviderCountryId())
	if config.GetDefaultTarget() == nil {
		config.DefaultTarget = &smsv1.SmsTarget{}
	}
	normalizeTarget(config.GetDefaultTarget())
	if config.GetCredentialSecret() == "" {
		existing, err := s.GetProviderConfig(ctx, config.GetProviderConfigId())
		if err == nil {
			config.CredentialSecret = existing.GetCredentialSecret()
		}
	}
	if config.GetEnabled() {
		if config.GetCredentialSecret() == "" {
			return nil, core.NewError(core.CodeValidationFailed, "credential_secret is required for enabled sms provider", false)
		}
	}
	if config.GetCapabilities() == nil {
		config.Capabilities = defaultProviderCapabilities(config.GetProviderKey())
	}
	config.Policy = providerPolicyToProto(providerPolicyFromConfig(config, defaultProviderPolicy(config.GetProviderKey())))
	config.Labels = normalizeLabels(config.GetLabels())
	return config, nil
}

func scanProviderConfig(row pgx.Row) (*smsinternalv1.SmsProviderConfig, error) {
	var config providerConfigRecord
	if err := row.Scan(
		&config.id, &config.providerKey, &config.displayName, &config.enabled, &config.apiEndpoint,
		&config.credentialSecret, &config.credentialSecretRef, &config.proxyRef, &config.httpProxy,
		&config.upstreamServiceKey, &config.providerCountryID,
		&config.defaultApplicationKey, &config.defaultCountryISO2, &config.defaultCountryCallingCode,
		&config.defaultMinPriceCurrency, &config.defaultMinPriceAmount,
		&config.defaultMaxPriceCurrency, &config.defaultMaxPriceAmount,
		&config.policyActivationTTLSeconds, &config.policyPollIntervalSeconds,
		&config.policyCancelAllowedAfterSeconds, &config.policyEarlyCancelRetryAfterSeconds,
		&config.policyCancelAllowedUntilSeconds, &config.capabilities, &config.labels,
		&config.createdAt, &config.updatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, core.NewError(core.CodeRouteNotFound, "sms provider config not found", false)
		}
		return nil, err
	}
	return config.toProto(), nil
}

type providerConfigRecord struct {
	id                                 string
	providerKey                        string
	displayName                        string
	enabled                            bool
	apiEndpoint                        string
	credentialSecret                   string
	credentialSecretRef                string
	proxyRef                           string
	httpProxy                          string
	upstreamServiceKey                 string
	providerCountryID                  string
	defaultApplicationKey              string
	defaultCountryISO2                 string
	defaultCountryCallingCode          string
	defaultMinPriceCurrency            string
	defaultMinPriceAmount              string
	defaultMaxPriceCurrency            string
	defaultMaxPriceAmount              string
	policyActivationTTLSeconds         int64
	policyPollIntervalSeconds          int64
	policyCancelAllowedAfterSeconds    int64
	policyEarlyCancelRetryAfterSeconds int64
	policyCancelAllowedUntilSeconds    int64
	capabilities                       []byte
	labels                             []byte
	createdAt                          time.Time
	updatedAt                          time.Time
}

func (r providerConfigRecord) toProto() *smsinternalv1.SmsProviderConfig {
	capabilities := &smsinternalv1.SmsProviderCapabilities{}
	_ = json.Unmarshal(r.capabilities, capabilities)
	labels := map[string]string{}
	_ = json.Unmarshal(r.labels, &labels)
	return &smsinternalv1.SmsProviderConfig{
		ProviderConfigId:    r.id,
		ProviderKey:         r.providerKey,
		DisplayName:         r.displayName,
		Enabled:             r.enabled,
		ApiEndpoint:         r.apiEndpoint,
		CredentialSecret:    r.credentialSecret,
		CredentialSecretRef: r.credentialSecretRef,
		ProxyRef:            r.proxyRef,
		HttpProxy:           r.httpProxy,
		UpstreamServiceKey:  r.upstreamServiceKey,
		ProviderCountryId:   r.providerCountryID,
		DefaultTarget: &smsv1.SmsTarget{
			ApplicationKey:     r.defaultApplicationKey,
			CountryIso2:        r.defaultCountryISO2,
			CountryCallingCode: r.defaultCountryCallingCode,
			MinPrice: &smsv1.DecimalMoney{
				CurrencyCode:  r.defaultMinPriceCurrency,
				AmountDecimal: r.defaultMinPriceAmount,
			},
			MaxPrice: &smsv1.DecimalMoney{
				CurrencyCode:  r.defaultMaxPriceCurrency,
				AmountDecimal: r.defaultMaxPriceAmount,
			},
		},
		Capabilities:        capabilities,
		Policy:              r.policyToProto(),
		Labels:              labels,
		CredentialSecretSet: r.credentialSecret != "",
		CreatedAt:           timestamppb.New(r.createdAt),
		UpdatedAt:           timestamppb.New(r.updatedAt),
	}
}

func (r providerConfigRecord) policyToProto() *smsinternalv1.SmsProviderPolicy {
	policy := defaultProviderPolicy(r.providerKey)
	if r.policyActivationTTLSeconds > 0 {
		policy.ActivationTTL = time.Duration(r.policyActivationTTLSeconds) * time.Second
	}
	if r.policyPollIntervalSeconds > 0 {
		policy.PollInterval = time.Duration(r.policyPollIntervalSeconds) * time.Second
	}
	if r.policyCancelAllowedAfterSeconds > 0 {
		policy.CancelAllowedAfter = time.Duration(r.policyCancelAllowedAfterSeconds) * time.Second
	}
	if r.policyEarlyCancelRetryAfterSeconds > 0 {
		policy.EarlyCancelRetryAfter = time.Duration(r.policyEarlyCancelRetryAfterSeconds) * time.Second
	}
	if r.policyCancelAllowedUntilSeconds > 0 {
		policy.CancelAllowedUntil = time.Duration(r.policyCancelAllowedUntilSeconds) * time.Second
	}
	return providerPolicyToProto(policy)
}

func providerConfigColumns() string {
	return `provider_config_id, provider_key, display_name, enabled, api_endpoint,
credential_secret, credential_secret_ref, proxy_ref, http_proxy,
upstream_service_key, provider_country_id,
default_application_key, default_country_iso2, default_country_calling_code,
default_min_price_currency, default_min_price_amount,
default_max_price_currency, default_max_price_amount,
policy_activation_ttl_seconds, policy_poll_interval_seconds,
policy_cancel_allowed_after_seconds, policy_early_cancel_retry_after_seconds,
policy_cancel_allowed_until_seconds, capabilities, labels,
created_at, updated_at`
}

func selectProviderConfigSQL() string {
	return ` RETURNING ` + providerConfigColumns()
}

func cloneProviderConfig(input *smsinternalv1.SmsProviderConfig) *smsinternalv1.SmsProviderConfig {
	if input == nil {
		return &smsinternalv1.SmsProviderConfig{}
	}
	return proto.Clone(input).(*smsinternalv1.SmsProviderConfig)
}

func normalizeTarget(target *smsv1.SmsTarget) {
	target.ApplicationKey = strings.TrimSpace(target.GetApplicationKey())
	target.CountryIso2 = strings.ToUpper(strings.TrimSpace(target.GetCountryIso2()))
	target.CountryCallingCode = strings.TrimPrefix(strings.TrimSpace(target.GetCountryCallingCode()), "+")
	if target.GetMaxPrice() != nil {
		target.MaxPrice.CurrencyCode = strings.ToUpper(strings.TrimSpace(target.GetMaxPrice().GetCurrencyCode()))
		target.MaxPrice.AmountDecimal = strings.TrimSpace(target.GetMaxPrice().GetAmountDecimal())
	}
	if target.GetMinPrice() != nil {
		target.MinPrice.CurrencyCode = strings.ToUpper(strings.TrimSpace(target.GetMinPrice().GetCurrencyCode()))
		target.MinPrice.AmountDecimal = strings.TrimSpace(target.GetMinPrice().GetAmountDecimal())
	}
}

func normalizeLabels(labels map[string]string) map[string]string {
	out := map[string]string{}
	for key, value := range labels {
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key != "" && value != "" {
			out[key] = value
		}
	}
	return out
}

func durationSeconds(value time.Duration) int64 {
	if value <= 0 {
		return 0
	}
	return int64(value / time.Second)
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

func targetMatchScore(configTarget *smsv1.SmsTarget, target core.Target) int {
	if configTarget == nil {
		return 0
	}
	score := 0
	if scoreField(configTarget.GetApplicationKey(), target.ApplicationKey, &score) < 0 {
		return -1
	}
	if scoreField(configTarget.GetCountryIso2(), target.CountryISO2, &score) < 0 {
		return -1
	}
	if scoreField(configTarget.GetCountryCallingCode(), target.CountryCallingCode, &score) < 0 {
		return -1
	}
	return score
}

func scoreField(expected, actual string, score *int) int {
	expected = strings.TrimSpace(expected)
	actual = strings.TrimSpace(actual)
	if expected == "" || actual == "" {
		return 0
	}
	if !strings.EqualFold(expected, actual) {
		return -1
	}
	*score = *score + 1
	return *score
}
