package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (p *ConfiguredProvider) GetBalance(ctx context.Context) (core.Money, error) {
	provider, err := p.providerForTarget(ctx, core.Target{})
	if err != nil {
		return core.Money{}, err
	}
	return provider.GetBalance(ctx)
}
