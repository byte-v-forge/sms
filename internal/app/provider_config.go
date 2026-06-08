package app

import (
	"context"
	"strings"
	"time"

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

func normalizeProviderKey(value string) string { return providerspi.NormalizeKey(value) }
