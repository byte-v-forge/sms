package spi

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
	"google.golang.org/protobuf/proto"
)

type Factory func(*smsinternalv1.SmsProviderConfig, *http.Client) (core.Provider, error)

type Plugin interface {
	Key() string
	DisplayName() string
	Capabilities() *smsinternalv1.SmsProviderCapabilities
	DefaultPolicy() core.ProviderPolicy
	NewProvider(*smsinternalv1.SmsProviderConfig, *http.Client) (core.Provider, error)
}

type Definition struct {
	ProviderKey   string
	DisplayName   string
	Capabilities  *smsinternalv1.SmsProviderCapabilities
	DefaultPolicy core.ProviderPolicy
	Factory       Factory
}

type definitionPlugin struct{ definition Definition }

type Registry struct {
	plugins map[string]Plugin
	keys    []string
}

func NewDefinition(definition Definition) Plugin {
	definition.ProviderKey = NormalizeKey(definition.ProviderKey)
	definition.DisplayName = strings.TrimSpace(definition.DisplayName)
	if definition.Capabilities == nil {
		definition.Capabilities = &smsinternalv1.SmsProviderCapabilities{}
	}
	return definitionPlugin{definition: definition}
}

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

func NormalizeKey(value string) string { return strings.ToLower(strings.TrimSpace(value)) }

func (p definitionPlugin) Key() string { return p.definition.ProviderKey }

func (p definitionPlugin) DisplayName() string { return p.definition.DisplayName }

func (p definitionPlugin) Capabilities() *smsinternalv1.SmsProviderCapabilities {
	return cloneCapabilities(p.definition.Capabilities)
}

func (p definitionPlugin) DefaultPolicy() core.ProviderPolicy { return p.definition.DefaultPolicy }

func (p definitionPlugin) NewProvider(config *smsinternalv1.SmsProviderConfig, client *http.Client) (core.Provider, error) {
	if p.definition.Factory == nil {
		return nil, core.NewError(core.CodeUnsupportedOperation, "sms provider factory is not configured", false)
	}
	return p.definition.Factory(config, client)
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

func BaseCapabilities(catalog bool) *smsinternalv1.SmsProviderCapabilities {
	return &smsinternalv1.SmsProviderCapabilities{
		SupportsBalance:         true,
		RequiresMarkMessageSent: true,
		SupportsAdditionalCode:  true,
		SupportsCatalog:         catalog,
		SupportsPriceLookup:     catalog,
	}
}

func cloneCapabilities(input *smsinternalv1.SmsProviderCapabilities) *smsinternalv1.SmsProviderCapabilities {
	if input == nil {
		return &smsinternalv1.SmsProviderCapabilities{}
	}
	return proto.Clone(input).(*smsinternalv1.SmsProviderCapabilities)
}
