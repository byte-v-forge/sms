package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (p *ConfiguredProvider) providerForTarget(ctx context.Context, target core.Target) (core.Provider, error) {
	config, err := p.configs.GetEnabledProviderConfig(ctx, p.key, target)
	if err != nil {
		return nil, err
	}
	return providerFromConfig(p.providers, config, p.timeout, p.defaultHTTPProxy)
}

func (p *ConfiguredProvider) providerForOrder(ctx context.Context, upstreamOrderID string) (core.Provider, error) {
	config, err := p.configs.GetProviderConfig(ctx, p.key)
	if err != nil {
		return nil, err
	}
	return providerFromConfig(p.providers, config, p.timeout, p.defaultHTTPProxy)
}
