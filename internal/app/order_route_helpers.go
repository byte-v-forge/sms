package app

import "github.com/byte-v-forge/sms/internal/core"

func withRouteTargetDefaults(target core.Target, route core.Route) core.Target {
	if target.ApplicationKey == "" {
		target.ApplicationKey = route.ApplicationKey
	}
	if target.CountryISO2 == "" {
		target.CountryISO2 = route.CountryISO2
	}
	if target.CountryCallingCode == "" {
		target.CountryCallingCode = route.CountryCallingCode
	}
	return target
}

func routeFromOrder(order core.Order) core.Route {
	return core.Route{
		ProviderKey:        order.ProviderKey,
		ApplicationKey:     order.Target.ApplicationKey,
		UpstreamServiceKey: order.Target.ApplicationKey,
		CountryISO2:        order.Target.CountryISO2,
		CountryCallingCode: order.Target.CountryCallingCode,
	}
}
