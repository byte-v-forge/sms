package app

import (
	"strings"
	"time"

	smsv1 "github.com/byte-v-forge/common-lib/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

func RouteFromPublicAcquireParams(params *smsv1.SmsNumberAcquireParams) core.Route {
	if params == nil {
		return core.Route{}
	}
	route := core.Route{
		ApplicationKey:     strings.TrimSpace(params.GetApplicationKey()),
		CountryISO2:        strings.ToUpper(strings.TrimSpace(params.GetCountryIso2())),
		CountryCallingCode: strings.TrimPrefix(strings.TrimSpace(params.GetCountryCallingCode()), "+"),
		FailurePolicy:      routeFailurePolicyFromProto(params.GetRouteFailurePolicy()),
	}
	switch value := params.GetProviderParams().(type) {
	case *smsv1.SmsNumberAcquireParams_FiveSim:
		five := value.FiveSim
		route.ProviderKey = "5sim"
		route.UpstreamServiceKey = strings.TrimSpace(five.GetProduct())
		route.ProviderCountryID = strings.TrimSpace(five.GetCountry())
		route.UpstreamProviderID = strings.TrimSpace(five.GetOperator())
	case *smsv1.SmsNumberAcquireParams_SmsBower:
		bower := value.SmsBower
		route.ProviderKey = "smsbower"
		route.UpstreamServiceKey = strings.TrimSpace(bower.GetService())
		route.ProviderCountryID = strings.TrimSpace(bower.GetCountry())
		route.UpstreamProviderID = strings.TrimSpace(bower.GetProviderId())
	case *smsv1.SmsNumberAcquireParams_HeroSms:
		hero := value.HeroSms
		route.ProviderKey = "herosms"
		route.UpstreamServiceKey = strings.TrimSpace(hero.GetService())
		route.ProviderCountryID = strings.TrimSpace(hero.GetCountry())
		route.UpstreamProviderID = strings.TrimSpace(hero.GetOperator())
	}
	return route
}

func PublicAcquireParamsFromRoute(route core.Route) *smsv1.SmsNumberAcquireParams {
	params := &smsv1.SmsNumberAcquireParams{
		ApplicationKey:     strings.TrimSpace(route.ApplicationKey),
		CountryIso2:        strings.ToUpper(strings.TrimSpace(route.CountryISO2)),
		CountryCallingCode: strings.TrimPrefix(strings.TrimSpace(route.CountryCallingCode), "+"),
	}
	if policy := protoRouteFailurePolicy(route.FailurePolicy); policy != nil {
		params.RouteFailurePolicy = policy
	}
	switch normalizeProviderKey(route.ProviderKey) {
	case "5sim":
		params.ProviderParams = &smsv1.SmsNumberAcquireParams_FiveSim{FiveSim: &smsv1.FiveSimAcquireParams{
			Product:  strings.TrimSpace(route.UpstreamServiceKey),
			Country:  strings.TrimSpace(route.ProviderCountryID),
			Operator: strings.TrimSpace(route.UpstreamProviderID),
		}}
	case "smsbower":
		params.ProviderParams = &smsv1.SmsNumberAcquireParams_SmsBower{SmsBower: &smsv1.SmsBowerAcquireParams{
			Service:    strings.TrimSpace(route.UpstreamServiceKey),
			Country:    strings.TrimSpace(route.ProviderCountryID),
			ProviderId: strings.TrimSpace(route.UpstreamProviderID),
		}}
	case "herosms":
		params.ProviderParams = &smsv1.SmsNumberAcquireParams_HeroSms{HeroSms: &smsv1.HeroSmsAcquireParams{
			Service:  strings.TrimSpace(route.UpstreamServiceKey),
			Country:  strings.TrimSpace(route.ProviderCountryID),
			Operator: strings.TrimSpace(route.UpstreamProviderID),
		}}
	}
	return params
}

func routeFailurePolicyFromRoutePolicy(policy *smsv1.SmsRoutePolicy) core.RouteFailurePolicy {
	if policy == nil {
		return core.RouteFailurePolicy{}
	}
	return routeFailurePolicyFromProto(policy.GetFailurePolicy())
}

func routeFailurePolicyFromProto(policy *smsv1.SmsRouteFailurePolicy) core.RouteFailurePolicy {
	if policy == nil {
		return core.RouteFailurePolicy{}
	}
	return core.RouteFailurePolicy{
		ScopeKey:         strings.TrimSpace(policy.GetScopeKey()),
		FailureThreshold: int(policy.GetFailureThreshold()),
		FailureWindow:    secondsDuration(policy.GetFailureWindowSeconds()),
		DisableTTL:       secondsDuration(policy.GetDisableTtlSeconds()),
	}
}

func protoRouteFailurePolicy(policy core.RouteFailurePolicy) *smsv1.SmsRouteFailurePolicy {
	if routeFailurePolicyIsZero(policy) {
		return nil
	}
	return &smsv1.SmsRouteFailurePolicy{
		ScopeKey:             strings.TrimSpace(policy.ScopeKey),
		FailureThreshold:     int32(policy.FailureThreshold),
		FailureWindowSeconds: int32(durationSeconds(policy.FailureWindow)),
		DisableTtlSeconds:    int32(durationSeconds(policy.DisableTTL)),
	}
}

func routeFailurePolicyIsZero(policy core.RouteFailurePolicy) bool {
	return strings.TrimSpace(policy.ScopeKey) == "" &&
		policy.FailureThreshold == 0 &&
		policy.FailureWindow <= 0 &&
		policy.DisableTTL <= 0
}

func secondsDuration(seconds int32) time.Duration {
	if seconds <= 0 {
		return 0
	}
	return time.Duration(seconds) * time.Second
}

func durationSeconds(duration time.Duration) int {
	if duration <= 0 {
		return 0
	}
	return int(duration.Round(time.Second) / time.Second)
}
