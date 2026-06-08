package app

import (
	"strings"

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
	route := core.Route{ProviderKey: normalizeProviderKey(ref.GetProviderKey())}
	target := ref.GetTarget()
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
	routeRef := ref.GetRouteRef()
	if routeRef != nil {
		route.UpstreamServiceKey = strings.TrimSpace(routeRef.GetUpstreamServiceKey())
		route.ProviderCountryID = strings.TrimSpace(routeRef.GetProviderCountryId())
		route.UpstreamProviderID = strings.TrimSpace(routeRef.GetUpstreamProviderId())
	}
	return route
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
		RouteRef: &smsv1.SmsOfferRouteRef{
			UpstreamServiceKey: strings.TrimSpace(route.UpstreamServiceKey),
			ProviderCountryId:  strings.TrimSpace(route.ProviderCountryID),
			UpstreamProviderId: strings.TrimSpace(route.UpstreamProviderID),
		},
	}
	if ref.GetProviderKey() == "" && ref.GetOfferId() == "" && targetIsZero(ref.GetTarget()) && offerRouteRefIsZero(ref.GetRouteRef()) {
		return nil
	}
	return ref
}
