package spi

import (
	"fmt"
	"sort"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

func NewRegistry(plugins ...Plugin) (*Registry, error) {
	registry := &Registry{plugins: make(map[string]Plugin, len(plugins))}
	for _, plugin := range plugins {
		if plugin == nil {
			return nil, fmt.Errorf("sms provider plugin is required")
		}
		key := NormalizeKey(plugin.Key())
		if key == "" {
			return nil, fmt.Errorf("sms provider plugin key is required")
		}
		if _, exists := registry.plugins[key]; exists {
			return nil, fmt.Errorf("duplicate sms provider plugin %q", key)
		}
		registry.plugins[key] = plugin
		registry.keys = append(registry.keys, key)
	}
	sort.Strings(registry.keys)
	return registry, nil
}

func EmptyRegistry() *Registry {
	return &Registry{plugins: map[string]Plugin{}, keys: []string{}}
}

func (r *Registry) Get(providerKey string) (Plugin, bool) {
	if r == nil {
		return nil, false
	}
	plugin, ok := r.plugins[NormalizeKey(providerKey)]
	return plugin, ok
}

func (r *Registry) Supports(providerKey string) bool {
	_, ok := r.Get(providerKey)
	return ok
}

func (r *Registry) Descriptors() []*smsinternalv1.SmsProviderPluginDescriptor {
	if r == nil {
		return nil
	}
	out := make([]*smsinternalv1.SmsProviderPluginDescriptor, 0, len(r.keys))
	for _, key := range r.keys {
		plugin := r.plugins[key]
		out = append(out, &smsinternalv1.SmsProviderPluginDescriptor{
			ProviderKey:  plugin.Key(),
			DisplayName:  plugin.DisplayName(),
			Capabilities: plugin.Capabilities(),
		})
	}
	return out
}

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
