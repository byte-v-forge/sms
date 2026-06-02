package app

import (
	"context"
	"net/http"
	"strings"
	"time"

	smsv1 "github.com/byte-v-forge/common-lib/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/common-lib/httpclient"
	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
	providerspi "github.com/byte-v-forge/sms/internal/providers/spi"
)

type ProviderConfigStore interface {
	UpsertProviderConfig(context.Context, *smsinternalv1.SmsProviderConfig) (*smsinternalv1.SmsProviderConfig, error)
	GetProviderConfig(context.Context, string) (*smsinternalv1.SmsProviderConfig, error)
	ListProviderConfigs(context.Context, bool, string) ([]*smsinternalv1.SmsProviderConfig, error)
	DeleteProviderConfig(context.Context, string) error
	GetEnabledProviderConfig(context.Context, string, core.Target) (*smsinternalv1.SmsProviderConfig, error)
}

type OrderListStore interface {
	List(context.Context, bool, int) ([]core.Order, error)
	ListCodes(context.Context, []string, int) ([]core.OrderCode, error)
}

type ConfiguredProvider struct {
	key              string
	providers        *providerspi.Registry
	configs          ProviderConfigStore
	timeout          time.Duration
	defaultHTTPProxy string
}

func NewConfiguredProvider(providers *providerspi.Registry, key string, configs ProviderConfigStore, timeout time.Duration, defaultHTTPProxy string) *ConfiguredProvider {
	return &ConfiguredProvider{key: normalizeProviderKey(key), providers: providers, configs: configs, timeout: timeout, defaultHTTPProxy: strings.TrimSpace(defaultHTTPProxy)}
}

func (p *ConfiguredProvider) Key() string { return p.key }
func (p *ConfiguredProvider) Policy() core.ProviderPolicy {
	return p.providers.DefaultPolicy(p.key, fallbackProviderPolicy())
}

func (p *ConfiguredProvider) AcquireNumber(ctx context.Context, request core.ProviderAcquireRequest) (core.ProviderOrder, error) {
	config, err := p.configs.GetEnabledProviderConfig(ctx, p.key, request.Target)
	if err != nil {
		return core.ProviderOrder{}, err
	}
	provider, err := providerFromConfig(p.providers, config, p.timeout, p.defaultHTTPProxy)
	if err != nil {
		return core.ProviderOrder{}, err
	}
	order, err := provider.AcquireNumber(ctx, request)
	if err == nil && order.UpstreamOrderID != "" {
		policy := provider.Policy().WithDefaults()
		if order.ExpiresAt.IsZero() && !order.AcquiredAt.IsZero() && policy.OrderTTL > 0 {
			order.ExpiresAt = order.AcquiredAt.Add(policy.OrderTTL)
		}
	}
	return order, err
}

func (p *ConfiguredProvider) GetStatus(ctx context.Context, upstreamOrderID string) (core.ProviderCodeResult, error) {
	provider, err := p.providerForOrder(ctx, upstreamOrderID)
	if err != nil {
		return core.ProviderCodeResult{}, err
	}
	return provider.GetStatus(ctx, upstreamOrderID)
}

func (p *ConfiguredProvider) SetStatus(ctx context.Context, upstreamOrderID string, action core.ProviderAction) error {
	provider, err := p.providerForOrder(ctx, upstreamOrderID)
	if err != nil {
		return err
	}
	return provider.SetStatus(ctx, upstreamOrderID, action)
}

func (p *ConfiguredProvider) BindOrderConfig(string) {}

func (p *ConfiguredProvider) PolicyForOrder(string) core.ProviderPolicy { return p.Policy() }
func (p *ConfiguredProvider) LoadPolicyForOrder(context.Context, string) core.ProviderPolicy {
	return p.Policy()
}

func (p *ConfiguredProvider) GetBalance(ctx context.Context) (core.Money, error) {
	config, err := p.configs.GetEnabledProviderConfig(ctx, p.key, core.Target{})
	if err != nil {
		return core.Money{}, err
	}
	provider, err := providerFromConfig(p.providers, config, p.timeout, p.defaultHTTPProxy)
	if err != nil {
		return core.Money{}, err
	}
	return provider.GetBalance(ctx)
}

func (p *ConfiguredProvider) providerForOrder(ctx context.Context, upstreamOrderID string) (core.Provider, error) {
	config, err := p.configs.GetProviderConfig(ctx, p.key)
	if err != nil {
		return nil, err
	}
	return providerFromConfig(p.providers, config, p.timeout, p.defaultHTTPProxy)
}

func providerFromConfig(providers *providerspi.Registry, config *smsinternalv1.SmsProviderConfig, timeout time.Duration, defaultHTTPProxy string) (core.Provider, error) {
	if config == nil {
		return nil, core.NewError(core.CodeRouteNotFound, "sms provider config not found", false)
	}
	client, err := httpClientFromConfig(timeout, defaultHTTPProxy)
	if err != nil {
		return nil, err
	}
	plugin, ok := providers.Get(config.GetProviderKey())
	if !ok {
		return nil, unsupportedSMSProvider(config.GetProviderKey())
	}
	return plugin.NewProvider(config, client)
}

func httpClientFromConfig(timeout time.Duration, defaultHTTPProxy string) (*http.Client, error) {
	if timeout <= 0 {
		timeout = 20 * time.Second
	}
	client, err := httpclient.New(timeout, strings.TrimSpace(defaultHTTPProxy))
	if err != nil {
		return nil, core.NewError(core.CodeValidationFailed, "invalid sms provider proxy", false)
	}
	return client, nil
}

func moneyFromProto(value *smsv1.DecimalMoney) core.Money {
	if value == nil {
		return core.Money{}
	}
	return core.Money{CurrencyCode: value.GetCurrencyCode(), AmountDecimal: value.GetAmountDecimal()}
}

func normalizeProviderKey(value string) string { return providerspi.NormalizeKey(value) }

func fallbackProviderPolicy() core.ProviderPolicy {
	return core.ProviderPolicy{OrderTTL: 20 * time.Minute, PollInterval: 5 * time.Second}
}
