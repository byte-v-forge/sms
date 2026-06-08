package app

import "github.com/byte-v-forge/sms/internal/core"

func (p *ConfiguredProvider) Policy() core.ProviderPolicy {
	return p.providers.DefaultPolicy(p.key, fallbackProviderPolicy())
}

func (p *ConfiguredProvider) BindOrderConfig(string) {}

func (p *ConfiguredProvider) PolicyForOrder(string) core.ProviderPolicy { return p.Policy() }
