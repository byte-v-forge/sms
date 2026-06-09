package app

import (
	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

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
