package app

import (
	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

func RouteFromPublicAcquireParams(params *smsv1.SmsNumberAcquireParams) core.Route {
	if params == nil {
		return core.Route{}
	}
	route := core.Route{
		ApplicationKey:     routeText(params.GetApplicationKey()),
		CountryISO2:        routeCountryISO2(params.GetCountryIso2()),
		CountryCallingCode: routeCallingCode(params.GetCountryCallingCode()),
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
			route.ApplicationKey = routeText(target.GetApplicationKey())
		}
		if route.CountryISO2 == "" {
			route.CountryISO2 = routeCountryISO2(target.GetCountryIso2())
		}
		if route.CountryCallingCode == "" {
			route.CountryCallingCode = routeCallingCode(target.GetCountryCallingCode())
		}
	}
	routeRef := ref.GetRouteRef()
	if routeRef != nil {
		route.UpstreamServiceKey = routeText(routeRef.GetUpstreamServiceKey())
		route.ProviderCountryID = routeText(routeRef.GetProviderCountryId())
		route.UpstreamProviderID = routeText(routeRef.GetUpstreamProviderId())
	}
	return route
}
