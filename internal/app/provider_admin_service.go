package app

import (
	"context"
	"strings"
	"time"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	providerspi "github.com/byte-v-forge/sms/internal/providers/spi"
)

type ProviderAdminService struct {
	configs          ProviderConfigStore
	providers        *providerspi.Registry
	orders           *OrderService
	orderDB          OrderListStore
	timeout          time.Duration
	defaultHTTPProxy string
	hot              HotStreamPublisher
}

func NewProviderAdminService(configs ProviderConfigStore, providers *providerspi.Registry, orders *OrderService, orderDB OrderListStore, timeout time.Duration, defaultHTTPProxy string, hot HotStreamPublisher) *ProviderAdminService {
	return &ProviderAdminService{
		configs:          configs,
		providers:        providers,
		orders:           orders,
		orderDB:          orderDB,
		timeout:          timeout,
		defaultHTTPProxy: strings.TrimSpace(defaultHTTPProxy),
		hot:              hot,
	}
}

func (s *ProviderAdminService) ListProviderPlugins(context.Context) ([]*smsinternalv1.SmsProviderPluginDescriptor, error) {
	return s.providers.Descriptors(), nil
}

func (s *ProviderAdminService) UpsertProviderConfig(ctx context.Context, config *smsinternalv1.SmsProviderConfig) (*smsinternalv1.SmsProviderConfig, error) {
	saved, err := s.configs.UpsertProviderConfig(ctx, config)
	if err != nil {
		return nil, err
	}
	s.publishProviderConfig(ctx, SMSProviderConfigUpdated, saved)
	return RedactProviderConfig(saved), nil
}

func (s *ProviderAdminService) GetProviderConfig(ctx context.Context, providerKey string) (*smsinternalv1.SmsProviderConfig, error) {
	config, err := s.configs.GetProviderConfig(ctx, providerKey)
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

func (s *ProviderAdminService) DeleteProviderConfig(ctx context.Context, providerKey string) error {
	if err := s.configs.DeleteProviderConfig(ctx, providerKey); err != nil {
		return err
	}
	s.publishResource(ctx, SMSProviderConfigDeleted, SMSProviderConfigResource, providerKey, map[string]string{"provider_key": providerKey})
	return nil
}
