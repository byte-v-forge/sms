package app

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/byte-v-forge/common-lib/eventbus"
	"github.com/byte-v-forge/common-lib/hotstream"
	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
	"google.golang.org/protobuf/proto"
)

type ProviderAdminService struct {
	configs          ProviderConfigStore
	orders           *OrderService
	orderDB          OrderListStore
	timeout          time.Duration
	defaultHTTPProxy string
	hot              HotStreamPublisher
}

func NewProviderAdminService(configs ProviderConfigStore, orders *OrderService, orderDB OrderListStore, timeout time.Duration, defaultHTTPProxy string, hot HotStreamPublisher) *ProviderAdminService {
	return &ProviderAdminService{
		configs:          configs,
		orders:           orders,
		orderDB:          orderDB,
		timeout:          timeout,
		defaultHTTPProxy: strings.TrimSpace(defaultHTTPProxy),
		hot:              hot,
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

func (s *ProviderAdminService) publishProviderConfig(ctx context.Context, eventType string, config *smsinternalv1.SmsProviderConfig) {
	if config == nil {
		return
	}
	s.publishResource(ctx, eventType, SMSProviderConfigResource, config.GetProviderKey(), map[string]string{
		"provider_key": config.GetProviderKey(),
		"enabled":      fmt.Sprintf("%t", config.GetEnabled()),
	})
}

func (s *ProviderAdminService) publishResource(ctx context.Context, eventType string, resourceType string, resourceID string, attrs map[string]string) {
	if s == nil || s.hot == nil || strings.TrimSpace(resourceID) == "" {
		return
	}
	now := time.Now()
	event := hotstream.NewEvent(hotstream.EventConfig{
		EventID:       eventbus.StableEventID("sms-hot-", eventType, resourceID, fmt.Sprintf("%d", now.UnixNano())),
		EventType:     eventType,
		SourceService: SMSHotStreamSource,
		ResourceType:  resourceType,
		ResourceID:    resourceID,
		OccurredAt:    now,
		CorrelationID: resourceID,
		Attributes:    attrs,
	})
	if err := s.hot.Publish(context.WithoutCancel(ctx), event); err != nil {
		log.Printf("publish sms config hotstream failed type=%s resource=%s: %v", eventType, resourceID, err)
	}
}

func (s *ProviderAdminService) GetProviderBalance(ctx context.Context, providerKey string) (core.Money, error) {
	config, err := s.configs.GetProviderConfig(ctx, providerKey)
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

func (s *ProviderAdminService) ListOrders(ctx context.Context, includeFinal bool, limit int) ([]core.Order, error) {
	if s.orderDB == nil {
		return nil, core.NewError(core.CodeUnsupportedOperation, "sms order list is not available", false)
	}
	if limit <= 0 {
		limit = 100
	}
	if limit > 500 {
		limit = 500
	}
	return s.orderDB.List(ctx, includeFinal, limit)
}

func (s *ProviderAdminService) CancelOrder(ctx context.Context, orderID string, requestID string) (core.Order, error) {
	if strings.TrimSpace(orderID) == "" {
		return core.Order{}, core.NewError(core.CodeValidationFailed, "order_id is required", false)
	}
	if strings.TrimSpace(requestID) == "" {
		requestID = RandomIDGenerator{}.NewID("req_")
	}
	return s.orders.CancelOrder(ctx, orderID, requestID)
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
