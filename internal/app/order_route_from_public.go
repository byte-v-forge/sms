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
	return mergeOfferRefRoute(route, RouteFromPublicOfferRef(params.GetOfferRef()))
}

func RouteFromPublicOfferRef(ref *smsv1.SmsOfferRef) core.Route {
	if ref == nil {
		return core.Route{}
	}
	route := core.Route{ProviderKey: normalizeProviderKey(ref.GetProviderKey())}
	target := ref.GetTarget()
	if target != nil {
		route = fillRouteTargetDefaults(route, routeFromPublicTarget(target))
	}
	routeRef := ref.GetRouteRef()
	if routeRef != nil {
		route.UpstreamServiceKey = routeText(routeRef.GetUpstreamServiceKey())
		route.ProviderCountryID = routeText(routeRef.GetProviderCountryId())
		route.UpstreamProviderID = routeText(routeRef.GetUpstreamProviderId())
	}
	return route
}

func mergeOfferRefRoute(route core.Route, refRoute core.Route) core.Route {
	if refRoute.ProviderKey == "" {
		return route
	}
	route.ProviderKey = refRoute.ProviderKey
	route.UpstreamServiceKey = refRoute.UpstreamServiceKey
	route.ProviderCountryID = refRoute.ProviderCountryID
	route.UpstreamProviderID = refRoute.UpstreamProviderID
	return fillRouteTargetDefaults(route, refRoute)
}

func routeFromPublicTarget(target *smsv1.SmsTarget) core.Route {
	return core.Route{
		ApplicationKey:     routeText(target.GetApplicationKey()),
		CountryISO2:        routeCountryISO2(target.GetCountryIso2()),
		CountryCallingCode: routeCallingCode(target.GetCountryCallingCode()),
	}
}

func fillRouteTargetDefaults(route core.Route, defaults core.Route) core.Route {
	if route.ApplicationKey == "" {
		route.ApplicationKey = defaults.ApplicationKey
	}
	if route.CountryISO2 == "" {
		route.CountryISO2 = defaults.CountryISO2
	}
	if route.CountryCallingCode == "" {
		route.CountryCallingCode = defaults.CountryCallingCode
	}
	return route
}
