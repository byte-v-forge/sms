package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

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
