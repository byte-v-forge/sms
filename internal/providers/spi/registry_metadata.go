package spi

import (
	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

func (r *Registry) Capabilities(providerKey string) *smsinternalv1.SmsProviderCapabilities {
	if plugin, ok := r.Get(providerKey); ok {
		return plugin.Capabilities()
	}
	return &smsinternalv1.SmsProviderCapabilities{}
}

func (r *Registry) DisplayName(providerKey string) string {
	if plugin, ok := r.Get(providerKey); ok {
		return plugin.DisplayName()
	}
	return NormalizeKey(providerKey)
}

func (r *Registry) DefaultPolicy(providerKey string, fallback core.ProviderPolicy) core.ProviderPolicy {
	if plugin, ok := r.Get(providerKey); ok {
		return plugin.DefaultPolicy()
	}
	return fallback
}
