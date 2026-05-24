package app

import (
	"context"
	"strings"
	"time"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
	"google.golang.org/protobuf/proto"
)

type ProviderAdminService struct {
	configs          ProviderConfigStore
	activations      *ActivationService
	activationDB     ActivationListStore
	timeout          time.Duration
	defaultHTTPProxy string
}

func NewProviderAdminService(configs ProviderConfigStore, activations *ActivationService, activationDB ActivationListStore, timeout time.Duration, defaultHTTPProxy string) *ProviderAdminService {
	return &ProviderAdminService{
		configs:          configs,
		activations:      activations,
		activationDB:     activationDB,
		timeout:          timeout,
		defaultHTTPProxy: strings.TrimSpace(defaultHTTPProxy),
	}
}

func (s *ProviderAdminService) ListProviderPlugins(context.Context) ([]*smsinternalv1.SmsProviderPluginDescriptor, error) {
	return listSMSProviderPluginDescriptors(), nil
}

func (s *ProviderAdminService) UpsertProviderConfig(ctx context.Context, config *smsinternalv1.SmsProviderConfig) (*smsinternalv1.SmsProviderConfig, error) {
	saved, err := s.configs.UpsertProviderConfig(ctx, config)
	if err != nil {
		return nil, err
	}
	return RedactProviderConfig(saved), nil
}

func (s *ProviderAdminService) GetProviderConfig(ctx context.Context, providerConfigID string) (*smsinternalv1.SmsProviderConfig, error) {
	config, err := s.configs.GetProviderConfig(ctx, providerConfigID)
	if err != nil {
		return nil, err
	}
	return RedactProviderConfig(config), nil
}

func (s *ProviderAdminService) ListProviderConfigs(ctx context.Context, includeDisabled bool, providerKey string) ([]*smsinternalv1.SmsProviderConfig, error) {
	configs, err := s.configs.ListProviderConfigs(ctx, includeDisabled, providerKey)
	if err != nil {
		return nil, err
	}
	for index, config := range configs {
		configs[index] = RedactProviderConfig(config)
	}
	return configs, nil
}

func (s *ProviderAdminService) DeleteProviderConfig(ctx context.Context, providerConfigID string) error {
	return s.configs.DeleteProviderConfig(ctx, providerConfigID)
}

func (s *ProviderAdminService) ListRouteOptions(ctx context.Context, providerConfigID, providerKey string) (*smsinternalv1.SmsProviderRouteOptions, error) {
	config, err := s.configForRouteOptions(ctx, providerConfigID, providerKey)
	if err != nil {
		return nil, err
	}
	if !config.GetEnabled() {
		return nil, core.NewError(core.CodeValidationFailed, "sms provider config is disabled", false)
	}
	provider, err := providerFromConfig(config, s.timeout, s.defaultHTTPProxy)
	if err != nil {
		return nil, err
	}
	return listRouteOptions(ctx, provider, config)
}

func (s *ProviderAdminService) UpsertRouteProfile(ctx context.Context, profile *smsinternalv1.SmsRouteProfile) (*smsinternalv1.SmsRouteProfile, error) {
	return s.configs.UpsertRouteProfile(ctx, profile)
}

func (s *ProviderAdminService) GetRouteProfile(ctx context.Context, profileKey string) (*smsinternalv1.SmsRouteProfile, error) {
	return s.configs.GetRouteProfile(ctx, profileKey)
}

func (s *ProviderAdminService) ListRouteProfiles(ctx context.Context, includeDisabled bool) ([]*smsinternalv1.SmsRouteProfile, error) {
	return s.configs.ListRouteProfiles(ctx, includeDisabled)
}

func (s *ProviderAdminService) DeleteRouteProfile(ctx context.Context, profileKey string) error {
	return s.configs.DeleteRouteProfile(ctx, profileKey)
}

func (s *ProviderAdminService) configForRouteOptions(ctx context.Context, providerConfigID, providerKey string) (*smsinternalv1.SmsProviderConfig, error) {
	if strings.TrimSpace(providerConfigID) != "" {
		return s.configs.GetProviderConfig(ctx, providerConfigID)
	}
	return s.configs.GetEnabledProviderConfig(ctx, providerKey, core.Target{})
}

func (s *ProviderAdminService) GetProviderBalance(ctx context.Context, providerConfigID string) (core.Money, error) {
	config, err := s.configs.GetProviderConfig(ctx, providerConfigID)
	if err != nil {
		return core.Money{}, err
	}
	if !config.GetEnabled() {
		return core.Money{}, core.NewError(core.CodeValidationFailed, "sms provider config is disabled", false)
	}
	provider, err := providerFromConfig(config, s.timeout, s.defaultHTTPProxy)
	if err != nil {
		return core.Money{}, err
	}
	return provider.GetBalance(ctx)
}

func (s *ProviderAdminService) ListActivations(ctx context.Context, includeFinal bool, limit int) ([]core.Activation, error) {
	if s.activationDB == nil {
		return nil, core.NewError(core.CodeUnsupportedOperation, "sms activation list is not available", false)
	}
	if limit <= 0 {
		limit = 100
	}
	if limit > 500 {
		limit = 500
	}
	return s.activationDB.List(ctx, includeFinal, limit)
}

func (s *ProviderAdminService) CancelActivation(ctx context.Context, activationID string, requestID string) (core.Activation, error) {
	if strings.TrimSpace(activationID) == "" {
		return core.Activation{}, core.NewError(core.CodeValidationFailed, "activation_id is required", false)
	}
	if strings.TrimSpace(requestID) == "" {
		requestID = RandomIDGenerator{}.NewID("req_")
	}
	return s.activations.CancelActivation(ctx, activationID, requestID)
}

func RedactProviderConfig(config *smsinternalv1.SmsProviderConfig) *smsinternalv1.SmsProviderConfig {
	if config == nil {
		return nil
	}
	redacted := proto.Clone(config).(*smsinternalv1.SmsProviderConfig)
	redacted.CredentialSecretSet = strings.TrimSpace(config.GetCredentialSecret()) != "" || config.GetCredentialSecretSet()
	redacted.CredentialSecret = ""
	return redacted
}
