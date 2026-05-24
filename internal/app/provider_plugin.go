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
	ConfigFields() []*smsinternalv1.SmsProviderConfigField
	RouteFields() []*smsinternalv1.SmsProviderRouteField
	RouteAdapter() routeProviderAdapter
	NewProvider(*smsinternalv1.SmsProviderConfig, *http.Client) (core.Provider, error)
}

type smsProviderPluginDefinition struct {
	key          string
	displayName  string
	capabilities *smsinternalv1.SmsProviderCapabilities
	policy       core.ProviderPolicy
	configFields []*smsinternalv1.SmsProviderConfigField
	routeFields  []*smsinternalv1.SmsProviderRouteField
	routeAdapter routeProviderAdapter
	newProvider  func(*smsinternalv1.SmsProviderConfig, *http.Client) (core.Provider, error)
}

func (p smsProviderPluginDefinition) Key() string         { return p.key }
func (p smsProviderPluginDefinition) DisplayName() string { return p.displayName }
func (p smsProviderPluginDefinition) Capabilities() *smsinternalv1.SmsProviderCapabilities {
	return proto.Clone(p.capabilities).(*smsinternalv1.SmsProviderCapabilities)
}
func (p smsProviderPluginDefinition) DefaultPolicy() core.ProviderPolicy { return p.policy }
func (p smsProviderPluginDefinition) ConfigFields() []*smsinternalv1.SmsProviderConfigField {
	return cloneConfigFields(p.configFields)
}
func (p smsProviderPluginDefinition) RouteFields() []*smsinternalv1.SmsProviderRouteField {
	return cloneRouteFields(p.routeFields)
}
func (p smsProviderPluginDefinition) RouteAdapter() routeProviderAdapter { return p.routeAdapter }
func (p smsProviderPluginDefinition) NewProvider(config *smsinternalv1.SmsProviderConfig, client *http.Client) (core.Provider, error) {
	return p.newProvider(config, client)
}

func smsProviderPlugins() []smsProviderPlugin {
	return []smsProviderPlugin{
		fiveSimPlugin(),
		heroSMSPlugin(),
		smsBowerPlugin(),
	}
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
			ConfigFields: plugin.ConfigFields(),
			RouteFields:  plugin.RouteFields(),
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
		policy: core.ProviderPolicy{ActivationTTL: 20 * time.Minute, PollInterval: 5 * time.Second},
		configFields: commonConfigFields(
			labelConfigField("currency_code", "Currency", smsinternalv1.SmsConfigFieldKind_SMS_CONFIG_FIELD_KIND_TEXT, false, true, "USD"),
		),
		routeFields: []*smsinternalv1.SmsProviderRouteField{
			routeField("upstream_service_key", "Product", smsinternalv1.SmsRouteFieldScope_SMS_ROUTE_FIELD_SCOPE_ROUTE, smsinternalv1.SmsRouteOptionSource_SMS_ROUTE_OPTION_SOURCE_SERVICES, false, ""),
			routeField("provider_country_id", "Country", smsinternalv1.SmsRouteFieldScope_SMS_ROUTE_FIELD_SCOPE_ROUTE, smsinternalv1.SmsRouteOptionSource_SMS_ROUTE_OPTION_SOURCE_COUNTRIES, false, ""),
			routeField("operator", "Operator", smsinternalv1.SmsRouteFieldScope_SMS_ROUTE_FIELD_SCOPE_OPTION, smsinternalv1.SmsRouteOptionSource_SMS_ROUTE_OPTION_SOURCE_OPERATORS, true, ""),
			routeField("amount_decimal", "Max Price", smsinternalv1.SmsRouteFieldScope_SMS_ROUTE_FIELD_SCOPE_MAX_PRICE, smsinternalv1.SmsRouteOptionSource_SMS_ROUTE_OPTION_SOURCE_UNSPECIFIED, true, ""),
			routeField("reuse", "Reuse", smsinternalv1.SmsRouteFieldScope_SMS_ROUTE_FIELD_SCOPE_OPTION, smsinternalv1.SmsRouteOptionSource_SMS_ROUTE_OPTION_SOURCE_UNSPECIFIED, true, ""),
			routeField("voice", "Voice", smsinternalv1.SmsRouteFieldScope_SMS_ROUTE_FIELD_SCOPE_OPTION, smsinternalv1.SmsRouteOptionSource_SMS_ROUTE_OPTION_SOURCE_UNSPECIFIED, true, ""),
			routeField("ref", "Ref", smsinternalv1.SmsRouteFieldScope_SMS_ROUTE_FIELD_SCOPE_OPTION, smsinternalv1.SmsRouteOptionSource_SMS_ROUTE_OPTION_SOURCE_UNSPECIFIED, true, ""),
		},
		routeAdapter: fivesim.RouteAdapter{},
		newProvider: func(config *smsinternalv1.SmsProviderConfig, client *http.Client) (core.Provider, error) {
			return fivesim.New(fivesim.Config{
				Endpoint:     config.GetApiEndpoint(),
				Token:        config.GetCredentialSecret(),
				CurrencyCode: firstLabel(config.GetLabels(), "currency", "currency_code"),
			}, client)
		},
	}
}

func heroSMSPlugin() smsProviderPlugin {
	return smsProviderPluginDefinition{
		key:          herosms.ProviderKey,
		displayName:  "HeroSMS",
		capabilities: baseCapabilities(false),
		policy:       core.ProviderPolicy{ActivationTTL: 20 * time.Minute, PollInterval: 5 * time.Second, CancelAllowedAfter: 2 * time.Minute},
		configFields: commonConfigFields(),
		routeFields: []*smsinternalv1.SmsProviderRouteField{
			routeField("upstream_service_key", "Service", smsinternalv1.SmsRouteFieldScope_SMS_ROUTE_FIELD_SCOPE_ROUTE, smsinternalv1.SmsRouteOptionSource_SMS_ROUTE_OPTION_SOURCE_SERVICES, false, ""),
			routeField("provider_country_id", "Country", smsinternalv1.SmsRouteFieldScope_SMS_ROUTE_FIELD_SCOPE_ROUTE, smsinternalv1.SmsRouteOptionSource_SMS_ROUTE_OPTION_SOURCE_COUNTRIES, false, ""),
			routeField("amount_decimal", "Max Price", smsinternalv1.SmsRouteFieldScope_SMS_ROUTE_FIELD_SCOPE_MAX_PRICE, smsinternalv1.SmsRouteOptionSource_SMS_ROUTE_OPTION_SOURCE_UNSPECIFIED, true, ""),
		},
		routeAdapter: herosms.RouteAdapter{},
		newProvider: func(config *smsinternalv1.SmsProviderConfig, client *http.Client) (core.Provider, error) {
			return herosms.New(herosms.Config{Endpoint: config.GetApiEndpoint(), APIKey: config.GetCredentialSecret()}, client)
		},
	}
}

func smsBowerPlugin() smsProviderPlugin {
	return smsProviderPluginDefinition{
		key:          smsbower.ProviderKey,
		displayName:  "SMSBower",
		capabilities: baseCapabilities(true),
		policy:       core.ProviderPolicy{ActivationTTL: 25 * time.Minute, PollInterval: 5 * time.Second, EarlyCancelRetryAfter: 2 * time.Minute},
		configFields: commonConfigFields(
			labelConfigField("ref", "Ref", smsinternalv1.SmsConfigFieldKind_SMS_CONFIG_FIELD_KIND_TEXT, false, true, ""),
			labelConfigField("user_id", "User ID", smsinternalv1.SmsConfigFieldKind_SMS_CONFIG_FIELD_KIND_TEXT, false, true, ""),
		),
		routeFields: []*smsinternalv1.SmsProviderRouteField{
			routeField("upstream_service_key", "Service", smsinternalv1.SmsRouteFieldScope_SMS_ROUTE_FIELD_SCOPE_ROUTE, smsinternalv1.SmsRouteOptionSource_SMS_ROUTE_OPTION_SOURCE_SERVICES, false, ""),
			routeField("provider_country_id", "Country", smsinternalv1.SmsRouteFieldScope_SMS_ROUTE_FIELD_SCOPE_ROUTE, smsinternalv1.SmsRouteOptionSource_SMS_ROUTE_OPTION_SOURCE_COUNTRIES, false, ""),
			routeField("amount_decimal", "Min Price", smsinternalv1.SmsRouteFieldScope_SMS_ROUTE_FIELD_SCOPE_MIN_PRICE, smsinternalv1.SmsRouteOptionSource_SMS_ROUTE_OPTION_SOURCE_UNSPECIFIED, true, ""),
			routeField("amount_decimal", "Max Price", smsinternalv1.SmsRouteFieldScope_SMS_ROUTE_FIELD_SCOPE_MAX_PRICE, smsinternalv1.SmsRouteOptionSource_SMS_ROUTE_OPTION_SOURCE_UNSPECIFIED, true, ""),
			routeField("ref", "Ref", smsinternalv1.SmsRouteFieldScope_SMS_ROUTE_FIELD_SCOPE_OPTION, smsinternalv1.SmsRouteOptionSource_SMS_ROUTE_OPTION_SOURCE_UNSPECIFIED, true, ""),
			routeField("include_provider_ids", "Include Providers", smsinternalv1.SmsRouteFieldScope_SMS_ROUTE_FIELD_SCOPE_OPTION, smsinternalv1.SmsRouteOptionSource_SMS_ROUTE_OPTION_SOURCE_UNSPECIFIED, true, ""),
			routeField("exclude_provider_ids", "Exclude Providers", smsinternalv1.SmsRouteFieldScope_SMS_ROUTE_FIELD_SCOPE_OPTION, smsinternalv1.SmsRouteOptionSource_SMS_ROUTE_OPTION_SOURCE_UNSPECIFIED, true, ""),
			routeField("phone_exception_prefixes", "Blocked Prefixes", smsinternalv1.SmsRouteFieldScope_SMS_ROUTE_FIELD_SCOPE_OPTION, smsinternalv1.SmsRouteOptionSource_SMS_ROUTE_OPTION_SOURCE_UNSPECIFIED, true, ""),
		},
		routeAdapter: smsbower.RouteAdapter{},
		newProvider: func(config *smsinternalv1.SmsProviderConfig, client *http.Client) (core.Provider, error) {
			return smsbower.New(smsbower.Config{
				Endpoint: config.GetApiEndpoint(),
				APIKey:   config.GetCredentialSecret(),
				Ref:      firstLabel(config.GetLabels(), "ref"),
				UserID:   firstLabel(config.GetLabels(), "userID", "user_id"),
			}, client)
		},
	}
}

func commonConfigFields(extra ...*smsinternalv1.SmsProviderConfigField) []*smsinternalv1.SmsProviderConfigField {
	fields := []*smsinternalv1.SmsProviderConfigField{
		configField(smsinternalv1.SmsConfigFieldTarget_SMS_CONFIG_FIELD_TARGET_CREDENTIAL_SECRET, "credential_secret", "API Key", smsinternalv1.SmsConfigFieldKind_SMS_CONFIG_FIELD_KIND_SECRET, true, false, ""),
		configField(smsinternalv1.SmsConfigFieldTarget_SMS_CONFIG_FIELD_TARGET_API_ENDPOINT, "api_endpoint", "API Endpoint", smsinternalv1.SmsConfigFieldKind_SMS_CONFIG_FIELD_KIND_URL, false, true, "provider default"),
		configField(smsinternalv1.SmsConfigFieldTarget_SMS_CONFIG_FIELD_TARGET_HTTP_PROXY, "http_proxy", "HTTP Proxy", smsinternalv1.SmsConfigFieldKind_SMS_CONFIG_FIELD_KIND_URL, false, true, "optional"),
	}
	return append(fields, extra...)
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

func labelConfigField(key, label string, kind smsinternalv1.SmsConfigFieldKind, required, advanced bool, placeholder string) *smsinternalv1.SmsProviderConfigField {
	return configField(smsinternalv1.SmsConfigFieldTarget_SMS_CONFIG_FIELD_TARGET_LABEL, key, label, kind, required, advanced, placeholder)
}

func configField(target smsinternalv1.SmsConfigFieldTarget, key, label string, kind smsinternalv1.SmsConfigFieldKind, required, advanced bool, placeholder string) *smsinternalv1.SmsProviderConfigField {
	return &smsinternalv1.SmsProviderConfigField{FieldKey: key, Label: label, Kind: kind, Required: required, Advanced: advanced, Placeholder: placeholder, Target: target}
}

func routeField(key, label string, scope smsinternalv1.SmsRouteFieldScope, options smsinternalv1.SmsRouteOptionSource, advanced bool, placeholder string) *smsinternalv1.SmsProviderRouteField {
	return &smsinternalv1.SmsProviderRouteField{FieldKey: key, Label: label, Scope: scope, OptionSource: options, Advanced: advanced, Placeholder: placeholder}
}

func cloneConfigFields(fields []*smsinternalv1.SmsProviderConfigField) []*smsinternalv1.SmsProviderConfigField {
	out := make([]*smsinternalv1.SmsProviderConfigField, 0, len(fields))
	for _, field := range fields {
		out = append(out, proto.Clone(field).(*smsinternalv1.SmsProviderConfigField))
	}
	return out
}

func cloneRouteFields(fields []*smsinternalv1.SmsProviderRouteField) []*smsinternalv1.SmsProviderRouteField {
	out := make([]*smsinternalv1.SmsProviderRouteField, 0, len(fields))
	for _, field := range fields {
		out = append(out, proto.Clone(field).(*smsinternalv1.SmsProviderRouteField))
	}
	return out
}

func unsupportedSMSProvider(providerKey string) error {
	return core.NewError(core.CodeUnsupportedOperation, fmt.Sprintf("unsupported sms provider %q", providerKey), false)
}
