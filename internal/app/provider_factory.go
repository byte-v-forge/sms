package app

import (
	"net/http"
	"strings"
	"time"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/httpclient"
	providerspi "github.com/byte-v-forge/sms/internal/providers/spi"
)

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

func fallbackProviderPolicy() core.ProviderPolicy {
	return core.ProviderPolicy{OrderTTL: 20 * time.Minute, PollInterval: 5 * time.Second}
}
