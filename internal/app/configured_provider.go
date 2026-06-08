package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

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
