package app

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
	"google.golang.org/protobuf/types/known/durationpb"
)

type ProviderConfigStore interface {
	UpsertProviderConfig(context.Context, *smsinternalv1.SmsProviderConfig) (*smsinternalv1.SmsProviderConfig, error)
	GetProviderConfig(context.Context, string) (*smsinternalv1.SmsProviderConfig, error)
	ListProviderConfigs(context.Context, bool, string) ([]*smsinternalv1.SmsProviderConfig, error)
	DeleteProviderConfig(context.Context, string) error
	UpsertRouteProfile(context.Context, *smsinternalv1.SmsRouteProfile) (*smsinternalv1.SmsRouteProfile, error)
	GetRouteProfile(context.Context, string) (*smsinternalv1.SmsRouteProfile, error)
	ListRouteProfiles(context.Context, bool) ([]*smsinternalv1.SmsRouteProfile, error)
	DeleteRouteProfile(context.Context, string) error
	GetEnabledProviderConfig(context.Context, string, core.Target) (*smsinternalv1.SmsProviderConfig, error)
}

type ActivationListStore interface {
	List(context.Context, bool, int) ([]core.Activation, error)
}

type ProviderConfigRouteResolver struct {
	configs ProviderConfigStore
}

func NewProviderConfigRouteResolver(configs ProviderConfigStore) *ProviderConfigRouteResolver {
	return &ProviderConfigRouteResolver{configs: configs}
}

func (r *ProviderConfigRouteResolver) Resolve(ctx context.Context, request core.RouteRequest) (core.Route, error) {
	if strings.TrimSpace(request.ProfileKey) != "" {
		return r.resolveProfile(ctx, request)
	}
	var config *smsinternalv1.SmsProviderConfig
	var err error
	if strings.TrimSpace(request.ProviderConfigID) != "" {
		config, err = r.configs.GetProviderConfig(ctx, request.ProviderConfigID)
		if err == nil {
			err = validateRequestedProviderConfig(config, request)
		}
	} else {
		config, err = r.configs.GetEnabledProviderConfig(ctx, request.ProviderKey, request.Target)
	}
	if err != nil {
		return core.Route{}, err
	}
	return routeFromProviderConfig(config, request.Target), nil
}

func (r *ProviderConfigRouteResolver) resolveProfile(ctx context.Context, request core.RouteRequest) (core.Route, error) {
	profile, err := r.configs.GetRouteProfile(ctx, request.ProfileKey)
	if err != nil {
		return core.Route{}, err
	}
	candidate, target, err := selectRouteCandidate(profile, request)
	if err != nil {
		return core.Route{}, err
	}
	config, err := r.configForCandidate(ctx, candidate, target)
	if err != nil {
		return core.Route{}, err
	}
	route := routeFromProviderConfig(config, target)
	route.ApplicationKey = target.ApplicationKey
	route.CountryISO2 = target.CountryISO2
	route.CountryCallingCode = target.CountryCallingCode
	applyRouteCandidate(candidate, &route)
	return route, nil
}

func (r *ProviderConfigRouteResolver) configForCandidate(ctx context.Context, candidate *smsinternalv1.SmsRouteCandidate, target core.Target) (*smsinternalv1.SmsProviderConfig, error) {
	if strings.TrimSpace(candidate.GetProviderConfigId()) != "" {
		config, err := r.configs.GetProviderConfig(ctx, candidate.GetProviderConfigId())
		if err != nil {
			return nil, err
		}
		if candidate.GetProviderKey() != "" && normalizeProviderKey(config.GetProviderKey()) != normalizeProviderKey(candidate.GetProviderKey()) {
			return nil, core.NewError(core.CodeRouteNotFound, "sms route provider config does not match provider", false)
		}
		return config, nil
	}
	return r.configs.GetEnabledProviderConfig(ctx, candidate.GetProviderKey(), target)
}

type ConfiguredProvider struct {
	key              string
	configs          ProviderConfigStore
	timeout          time.Duration
	defaultHTTPProxy string
	activation       sync.Map
	policy           sync.Map
}

func NewConfiguredProvider(key string, configs ProviderConfigStore, timeout time.Duration, defaultHTTPProxy string) *ConfiguredProvider {
	return &ConfiguredProvider{
		key:              normalizeProviderKey(key),
		configs:          configs,
		timeout:          timeout,
		defaultHTTPProxy: strings.TrimSpace(defaultHTTPProxy),
	}
}

func (p *ConfiguredProvider) Key() string {
	return p.key
}

func (p *ConfiguredProvider) Policy() core.ProviderPolicy {
	return defaultProviderPolicy(p.key)
}

func (p *ConfiguredProvider) AcquireNumber(ctx context.Context, request core.ProviderAcquireRequest) (core.ProviderActivation, error) {
	config, err := p.configForRoute(ctx, request.Route, request.Target)
	if err != nil {
		return core.ProviderActivation{}, err
	}
	provider, err := providerFromConfig(config, p.timeout, p.defaultHTTPProxy)
	if err != nil {
		return core.ProviderActivation{}, err
	}
	policy := providerPolicyFromConfig(config, provider.Policy())
	activation, err := provider.AcquireNumber(ctx, request)
	if err == nil && activation.UpstreamActivationID != "" {
		p.activation.Store(activation.UpstreamActivationID, config.GetProviderConfigId())
		p.policy.Store(activation.UpstreamActivationID, policy)
		if activation.ExpiresAt.IsZero() && !activation.AcquiredAt.IsZero() && policy.ActivationTTL > 0 {
			activation.ExpiresAt = activation.AcquiredAt.Add(policy.ActivationTTL)
		}
	}
	return activation, err
}

func (p *ConfiguredProvider) GetStatus(ctx context.Context, upstreamActivationID string) (core.ProviderCodeResult, error) {
	provider, err := p.providerForActivation(ctx, upstreamActivationID)
	if err != nil {
		return core.ProviderCodeResult{}, err
	}
	return provider.GetStatus(ctx, upstreamActivationID)
}

func (p *ConfiguredProvider) SetStatus(ctx context.Context, upstreamActivationID string, action core.ProviderAction) error {
	provider, err := p.providerForActivation(ctx, upstreamActivationID)
	if err != nil {
		return err
	}
	return provider.SetStatus(ctx, upstreamActivationID, action)
}

func (p *ConfiguredProvider) BindActivationConfig(upstreamActivationID string, providerConfigID string) {
	upstreamActivationID = strings.TrimSpace(upstreamActivationID)
	providerConfigID = strings.TrimSpace(providerConfigID)
	if upstreamActivationID == "" || providerConfigID == "" {
		return
	}
	p.activation.Store(upstreamActivationID, providerConfigID)
}

func (p *ConfiguredProvider) PolicyForActivation(upstreamActivationID string) core.ProviderPolicy {
	upstreamActivationID = strings.TrimSpace(upstreamActivationID)
	if upstreamActivationID != "" {
		if value, ok := p.policy.Load(upstreamActivationID); ok {
			if policy, ok := value.(core.ProviderPolicy); ok {
				return policy
			}
		}
	}
	return p.Policy()
}

func (p *ConfiguredProvider) LoadPolicyForActivation(ctx context.Context, upstreamActivationID string, providerConfigID string) core.ProviderPolicy {
	upstreamActivationID = strings.TrimSpace(upstreamActivationID)
	providerConfigID = strings.TrimSpace(providerConfigID)
	if upstreamActivationID != "" {
		if value, ok := p.policy.Load(upstreamActivationID); ok {
			if policy, ok := value.(core.ProviderPolicy); ok {
				return policy
			}
		}
	}
	if providerConfigID == "" {
		return p.Policy()
	}
	config, err := p.configs.GetProviderConfig(ctx, providerConfigID)
	if err != nil {
		return p.Policy()
	}
	policy := providerPolicyFromConfig(config, defaultProviderPolicy(config.GetProviderKey()))
	if upstreamActivationID != "" {
		p.activation.Store(upstreamActivationID, providerConfigID)
		p.policy.Store(upstreamActivationID, policy)
	}
	return policy
}

func (p *ConfiguredProvider) GetBalance(ctx context.Context) (core.Money, error) {
	config, err := p.configs.GetEnabledProviderConfig(ctx, p.key, core.Target{})
	if err != nil {
		return core.Money{}, err
	}
	provider, err := providerFromConfig(config, p.timeout, p.defaultHTTPProxy)
	if err != nil {
		return core.Money{}, err
	}
	return provider.GetBalance(ctx)
}

func (p *ConfiguredProvider) configForRoute(ctx context.Context, route core.Route, target core.Target) (*smsinternalv1.SmsProviderConfig, error) {
	configID := strings.TrimSpace(route.ProviderOptions["provider_config_id"])
	if configID != "" {
		return p.configs.GetProviderConfig(ctx, configID)
	}
	return p.configs.GetEnabledProviderConfig(ctx, p.key, target)
}

func (p *ConfiguredProvider) providerForActivation(ctx context.Context, upstreamActivationID string) (core.Provider, error) {
	if value, ok := p.activation.Load(upstreamActivationID); ok {
		if configID, ok := value.(string); ok && strings.TrimSpace(configID) != "" {
			config, err := p.configs.GetProviderConfig(ctx, configID)
			if err != nil {
				return nil, err
			}
			provider, err := providerFromConfig(config, p.timeout, p.defaultHTTPProxy)
			if err == nil {
				p.policy.Store(upstreamActivationID, providerPolicyFromConfig(config, provider.Policy()))
			}
			return provider, err
		}
	}
	config, err := p.configs.GetEnabledProviderConfig(ctx, p.key, core.Target{})
	if err != nil {
		return nil, err
	}
	provider, err := providerFromConfig(config, p.timeout, p.defaultHTTPProxy)
	if err == nil {
		p.policy.Store(upstreamActivationID, providerPolicyFromConfig(config, provider.Policy()))
	}
	return provider, err
}

func providerFromConfig(config *smsinternalv1.SmsProviderConfig, timeout time.Duration, defaultHTTPProxy string) (core.Provider, error) {
	if config == nil {
		return nil, core.NewError(core.CodeRouteNotFound, "sms provider config not found", false)
	}
	client, err := httpClientFromConfig(config, timeout, defaultHTTPProxy)
	if err != nil {
		return nil, err
	}
	plugin, ok := smsProviderPluginByKey(config.GetProviderKey())
	if !ok {
		return nil, unsupportedSMSProvider(config.GetProviderKey())
	}
	return plugin.NewProvider(config, client)
}

func routeFromProviderConfig(config *smsinternalv1.SmsProviderConfig, requestTarget core.Target) core.Route {
	defaultTarget := targetFromProto(config.GetDefaultTarget())
	target := mergeRouteTarget(requestTarget, defaultTarget)
	options := cloneMap(config.GetLabels())
	if options == nil {
		options = map[string]string{}
	}
	options["provider_config_id"] = config.GetProviderConfigId()
	return core.Route{
		ProviderKey:        normalizeProviderKey(config.GetProviderKey()),
		ApplicationKey:     target.ApplicationKey,
		UpstreamServiceKey: firstNonEmptyString(config.GetUpstreamServiceKey(), target.ApplicationKey),
		CountryISO2:        target.CountryISO2,
		CountryCallingCode: target.CountryCallingCode,
		ProviderCountryID:  config.GetProviderCountryId(),
		MinPrice:           target.MinPrice,
		MaxPrice:           target.MaxPrice,
		ProviderOptions:    options,
	}
}

func httpClientFromConfig(config *smsinternalv1.SmsProviderConfig, timeout time.Duration, defaultHTTPProxy string) (*http.Client, error) {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	if proxy := firstNonEmptyString(config.GetHttpProxy(), defaultHTTPProxy); proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			return nil, core.NewError(core.CodeValidationFailed, "invalid sms provider proxy", false)
		}
		transport.Proxy = http.ProxyURL(proxyURL)
	}
	if timeout <= 0 {
		timeout = 20 * time.Second
	}
	return &http.Client{Transport: transport, Timeout: timeout}, nil
}

func moneyFromProto(value *smsv1.DecimalMoney) core.Money {
	if value == nil {
		return core.Money{}
	}
	return core.Money{CurrencyCode: value.GetCurrencyCode(), AmountDecimal: value.GetAmountDecimal()}
}

func targetFromProto(value *smsv1.SmsTarget) core.Target {
	if value == nil {
		return core.Target{}
	}
	return core.Target{
		ApplicationKey:     strings.TrimSpace(value.GetApplicationKey()),
		CountryISO2:        strings.ToUpper(strings.TrimSpace(value.GetCountryIso2())),
		CountryCallingCode: strings.TrimPrefix(strings.TrimSpace(value.GetCountryCallingCode()), "+"),
		MinPrice:           moneyFromProto(value.GetMinPrice()),
		MaxPrice:           moneyFromProto(value.GetMaxPrice()),
	}
}

func mergeRouteTarget(primary core.Target, fallback core.Target) core.Target {
	if primary.ApplicationKey == "" {
		primary.ApplicationKey = fallback.ApplicationKey
	}
	if primary.CountryISO2 == "" {
		primary.CountryISO2 = fallback.CountryISO2
	}
	if primary.CountryCallingCode == "" {
		primary.CountryCallingCode = fallback.CountryCallingCode
	}
	if primary.MinPrice.AmountDecimal == "" && primary.MinPrice.CurrencyCode == "" {
		primary.MinPrice = fallback.MinPrice
	}
	if primary.MaxPrice.AmountDecimal == "" && primary.MaxPrice.CurrencyCode == "" {
		primary.MaxPrice = fallback.MaxPrice
	}
	return primary
}

func firstLabel(labels map[string]string, keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(labels[key]); value != "" {
			return value
		}
	}
	return ""
}

func normalizeProviderKey(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			return value
		}
	}
	return ""
}

func defaultProviderPolicy(providerKey string) core.ProviderPolicy {
	if plugin, ok := smsProviderPluginByKey(providerKey); ok {
		return plugin.DefaultPolicy()
	}
	return core.ProviderPolicy{ActivationTTL: 20 * time.Minute, PollInterval: 5 * time.Second}
}

func providerPolicyFromConfig(config *smsinternalv1.SmsProviderConfig, fallback core.ProviderPolicy) core.ProviderPolicy {
	policy := fallback.WithDefaults()
	if config == nil || config.GetPolicy() == nil {
		return policy
	}
	if value := protoDuration(config.GetPolicy().GetActivationTtl()); value > 0 {
		policy.ActivationTTL = value
	}
	if value := protoDuration(config.GetPolicy().GetPollInterval()); value > 0 {
		policy.PollInterval = value
	}
	if value := protoDuration(config.GetPolicy().GetCancelAllowedAfter()); value > 0 {
		policy.CancelAllowedAfter = value
	}
	if value := protoDuration(config.GetPolicy().GetEarlyCancelRetryAfter()); value > 0 {
		policy.EarlyCancelRetryAfter = value
	}
	if value := protoDuration(config.GetPolicy().GetCancelAllowedUntil()); value > 0 {
		policy.CancelAllowedUntil = value
	}
	return policy
}

func providerPolicyToProto(policy core.ProviderPolicy) *smsinternalv1.SmsProviderPolicy {
	policy = policy.WithDefaults()
	return &smsinternalv1.SmsProviderPolicy{
		ActivationTtl:         durationpb.New(policy.ActivationTTL),
		PollInterval:          durationpb.New(policy.PollInterval),
		CancelAllowedAfter:    durationOrNil(policy.CancelAllowedAfter),
		EarlyCancelRetryAfter: durationOrNil(policy.EarlyCancelRetryAfter),
		CancelAllowedUntil:    durationOrNil(policy.CancelAllowedUntil),
	}
}

func protoDuration(value *durationpb.Duration) time.Duration {
	if value == nil {
		return 0
	}
	return value.AsDuration()
}

func durationOrNil(value time.Duration) *durationpb.Duration {
	if value <= 0 {
		return nil
	}
	return durationpb.New(value)
}

func validateRequestedProviderConfig(config *smsinternalv1.SmsProviderConfig, request core.RouteRequest) error {
	if config == nil || !config.GetEnabled() {
		return core.NewError(core.CodeRouteNotFound, "sms provider config is disabled or missing", false)
	}
	if providerKey := normalizeProviderKey(request.ProviderKey); providerKey != "" && providerKey != normalizeProviderKey(config.GetProviderKey()) {
		return core.NewError(core.CodeRouteNotFound, "sms provider config does not match requested provider", false)
	}
	if targetMatchScore(config.GetDefaultTarget(), request.Target) < 0 {
		return core.NewError(core.CodeRouteNotFound, "sms provider config does not match requested target", false)
	}
	return nil
}
