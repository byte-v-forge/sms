package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (p *ConfiguredProvider) LoadPolicyForOrder(context.Context, string) core.ProviderPolicy {
	return p.Policy()
}
