package app

import (
	"strings"
	"time"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
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
		MinAvailableCount:  int(params.GetMinAvailableCount()),
		MinPrice:           moneyFromProto(params.GetMinPrice()),
		MaxPrice:           moneyFromProto(params.GetMaxPrice()),
		FailurePolicy:      routeFailurePolicyFromProto(params.GetRouteFailurePolicy()),
	}
	if refRoute := RouteFromPublicOfferRef(params.GetOfferRef()); refRoute.ProviderKey != "" {
		route.ProviderKey = refRoute.ProviderKey
		route.UpstreamServiceKey = refRoute.UpstreamServiceKey
		route.ProviderCountryID = refRoute.ProviderCountryID
		route.UpstreamProviderID = refRoute.UpstreamProviderID
		if route.ApplicationKey == "" {
			route.ApplicationKey = refRoute.ApplicationKey
		}
		if route.CountryISO2 == "" {
			route.CountryISO2 = refRoute.CountryISO2
		}
		if route.CountryCallingCode == "" {
			route.CountryCallingCode = refRoute.CountryCallingCode
		}
	}
	return route
}

func RouteFromPublicOfferRef(ref *smsv1.SmsOfferRef) core.Route {
	if ref == nil {
		return core.Route{}
	}
	route := routeFromOfferID(ref.GetOfferId())
	target := ref.GetTarget()
	if route.ProviderKey == "" {
		route.ProviderKey = normalizeProviderKey(ref.GetProviderKey())
	}
	if target != nil {
		if route.ApplicationKey == "" {
			route.ApplicationKey = strings.TrimSpace(target.GetApplicationKey())
		}
		if route.CountryISO2 == "" {
			route.CountryISO2 = strings.ToUpper(strings.TrimSpace(target.GetCountryIso2()))
		}
		if route.CountryCallingCode == "" {
			route.CountryCallingCode = strings.TrimPrefix(strings.TrimSpace(target.GetCountryCallingCode()), "+")
		}
	}
	return route
}

func routeFromOfferID(offerID string) core.Route {
	parts := strings.Split(strings.TrimSpace(offerID), "|")
	if len(parts) != 7 {
		return core.Route{}
	}
	return core.Route{
		ProviderKey:        normalizeProviderKey(parts[0]),
		ApplicationKey:     strings.TrimSpace(parts[1]),
		CountryISO2:        strings.ToUpper(strings.TrimSpace(parts[2])),
		CountryCallingCode: strings.TrimPrefix(strings.TrimSpace(parts[3]), "+"),
		UpstreamServiceKey: strings.TrimSpace(parts[4]),
		ProviderCountryID:  strings.TrimSpace(parts[5]),
		UpstreamProviderID: strings.TrimSpace(parts[6]),
	}
}

func PublicAcquireParamsFromRoute(route core.Route) *smsv1.SmsNumberAcquireParams {
	params := &smsv1.SmsNumberAcquireParams{
		OfferRef:           PublicOfferRefFromRoute(route),
		ApplicationKey:     strings.TrimSpace(route.ApplicationKey),
		CountryIso2:        strings.ToUpper(strings.TrimSpace(route.CountryISO2)),
		CountryCallingCode: strings.TrimPrefix(strings.TrimSpace(route.CountryCallingCode), "+"),
		MinAvailableCount:  int32(route.MinAvailableCount),
	}
	if moneyIsSet(route.MinPrice) {
		params.MinPrice = PublicMoney(route.MinPrice)
	}
	if moneyIsSet(route.MaxPrice) {
		params.MaxPrice = PublicMoney(route.MaxPrice)
	}
	if policy := protoRouteFailurePolicy(route.FailurePolicy); policy != nil {
		params.RouteFailurePolicy = policy
	}
	return params
}

func PublicOfferRefFromRoute(route core.Route) *smsv1.SmsOfferRef {
	ref := &smsv1.SmsOfferRef{
		OfferId:     publicOfferID(route),
		ProviderKey: normalizeProviderKey(route.ProviderKey),
		Target: &smsv1.SmsTarget{
			ApplicationKey:     strings.TrimSpace(route.ApplicationKey),
			CountryIso2:        strings.ToUpper(strings.TrimSpace(route.CountryISO2)),
			CountryCallingCode: strings.TrimPrefix(strings.TrimSpace(route.CountryCallingCode), "+"),
		},
	}
	if ref.GetProviderKey() == "" && ref.GetOfferId() == "" && ref.GetTarget().GetApplicationKey() == "" && ref.GetTarget().GetCountryIso2() == "" && ref.GetTarget().GetCountryCallingCode() == "" {
		return nil
	}
	return ref
}

func publicOfferID(route core.Route) string {
	values := []string{
		normalizeProviderKey(route.ProviderKey),
		strings.TrimSpace(route.ApplicationKey),
		strings.ToUpper(strings.TrimSpace(route.CountryISO2)),
		strings.TrimPrefix(strings.TrimSpace(route.CountryCallingCode), "+"),
		strings.TrimSpace(route.UpstreamServiceKey),
		strings.TrimSpace(route.ProviderCountryID),
		strings.TrimSpace(route.UpstreamProviderID),
	}
	if strings.Join(values, "") == "" {
		return ""
	}
	return strings.Join(values, "|")
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

func moneyIsSet(money core.Money) bool {
	return strings.TrimSpace(money.AmountDecimal) != "" || strings.TrimSpace(money.CurrencyCode) != ""
}
