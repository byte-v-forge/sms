package app

import (
	"fmt"
	"net/http"
	"time"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/providers/fivesim"
	"github.com/byte-v-forge/sms/internal/providers/herosms"
	"github.com/byte-v-forge/sms/internal/providers/smsbower"
	"google.golang.org/protobuf/proto"
)

type smsProviderPlugin interface {
	Key() string
	DisplayName() string
	Capabilities() *smsinternalv1.SmsProviderCapabilities
	DefaultPolicy() core.ProviderPolicy
	NewProvider(*smsinternalv1.SmsProviderConfig, *http.Client) (core.Provider, error)
}

type smsProviderPluginDefinition struct {
	key          string
	displayName  string
	capabilities *smsinternalv1.SmsProviderCapabilities
	policy       core.ProviderPolicy
	newProvider  func(*smsinternalv1.SmsProviderConfig, *http.Client) (core.Provider, error)
}

func (p smsProviderPluginDefinition) Key() string         { return p.key }
func (p smsProviderPluginDefinition) DisplayName() string { return p.displayName }
func (p smsProviderPluginDefinition) Capabilities() *smsinternalv1.SmsProviderCapabilities {
	return proto.Clone(p.capabilities).(*smsinternalv1.SmsProviderCapabilities)
}
func (p smsProviderPluginDefinition) DefaultPolicy() core.ProviderPolicy { return p.policy }
func (p smsProviderPluginDefinition) NewProvider(config *smsinternalv1.SmsProviderConfig, client *http.Client) (core.Provider, error) {
	return p.newProvider(config, client)
}

func smsProviderPlugins() []smsProviderPlugin {
	return []smsProviderPlugin{fiveSimPlugin(), heroSMSPlugin(), smsBowerPlugin()}
}

func NewConfiguredProviders(configs ProviderConfigStore, timeout time.Duration, defaultHTTPProxy string) []core.Provider {
	plugins := smsProviderPlugins()
	providers := make([]core.Provider, 0, len(plugins))
	for _, plugin := range plugins {
		providers = append(providers, NewConfiguredProvider(plugin.Key(), configs, timeout, defaultHTTPProxy))
	}
	return providers
}

func smsProviderPluginByKey(providerKey string) (smsProviderPlugin, bool) {
	key := normalizeProviderKey(providerKey)
	for _, plugin := range smsProviderPlugins() {
		if plugin.Key() == key {
			return plugin, true
		}
	}
	return nil, false
}

func listSMSProviderPluginDescriptors() []*smsinternalv1.SmsProviderPluginDescriptor {
	plugins := smsProviderPlugins()
	descriptors := make([]*smsinternalv1.SmsProviderPluginDescriptor, 0, len(plugins))
	for _, plugin := range plugins {
		descriptors = append(descriptors, &smsinternalv1.SmsProviderPluginDescriptor{
			ProviderKey:  plugin.Key(),
			DisplayName:  plugin.DisplayName(),
			Capabilities: plugin.Capabilities(),
		})
	}
	return descriptors
}

func fiveSimPlugin() smsProviderPlugin {
	return smsProviderPluginDefinition{
		key:         fivesim.ProviderKey,
		displayName: "5sim",
		capabilities: &smsinternalv1.SmsProviderCapabilities{
			SupportsBalance:         true,
			RequiresMarkMessageSent: true,
			SupportsAdditionalCode:  true,
			SupportsCatalog:         true,
			SupportsPriceLookup:     true,
		},
		policy: core.ProviderPolicy{OrderTTL: 20 * time.Minute, PollInterval: 5 * time.Second},
		newProvider: func(config *smsinternalv1.SmsProviderConfig, client *http.Client) (core.Provider, error) {
			return fivesim.New(fivesim.Config{Token: config.GetCredentialSecret()}, client)
		},
	}
}

func heroSMSPlugin() smsProviderPlugin {
	return smsProviderPluginDefinition{
		key:          herosms.ProviderKey,
		displayName:  "HeroSMS",
		capabilities: baseCapabilities(false),
		policy:       core.ProviderPolicy{OrderTTL: 20 * time.Minute, PollInterval: 5 * time.Second, CancelAllowedAfter: 2 * time.Minute},
		newProvider: func(config *smsinternalv1.SmsProviderConfig, client *http.Client) (core.Provider, error) {
			return herosms.New(herosms.Config{APIKey: config.GetCredentialSecret()}, client)
		},
	}
}

func smsBowerPlugin() smsProviderPlugin {
	return smsProviderPluginDefinition{
		key:          smsbower.ProviderKey,
		displayName:  "SMSBower",
		capabilities: baseCapabilities(true),
		policy:       core.ProviderPolicy{OrderTTL: 25 * time.Minute, PollInterval: 5 * time.Second, EarlyCancelRetryAfter: 2 * time.Minute},
		newProvider: func(config *smsinternalv1.SmsProviderConfig, client *http.Client) (core.Provider, error) {
			return smsbower.New(smsbower.Config{APIKey: config.GetCredentialSecret()}, client)
		},
	}
}

func baseCapabilities(catalog bool) *smsinternalv1.SmsProviderCapabilities {
	return &smsinternalv1.SmsProviderCapabilities{
		SupportsBalance:         true,
		RequiresMarkMessageSent: true,
		SupportsAdditionalCode:  true,
		SupportsCatalog:         catalog,
		SupportsPriceLookup:     catalog,
	}
}

func unsupportedSMSProvider(providerKey string) error {
	return core.NewError(core.CodeUnsupportedOperation, fmt.Sprintf("unsupported sms provider %q", providerKey), false)
}
