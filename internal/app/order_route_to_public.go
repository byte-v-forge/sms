package app

import (
	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

func PublicAcquireParamsFromRoute(route core.Route) *smsv1.SmsNumberAcquireParams {
	params := &smsv1.SmsNumberAcquireParams{
		OfferRef:           PublicOfferRefFromRoute(route),
		ApplicationKey:     routeText(route.ApplicationKey),
		CountryIso2:        routeCountryISO2(route.CountryISO2),
		CountryCallingCode: routeCallingCode(route.CountryCallingCode),
		MinAvailableCount:  int32(route.MinAvailableCount),
	}
	applyPublicAcquirePriceParams(params, route)
	if policy := protoRouteFailurePolicy(route.FailurePolicy); policy != nil {
		params.RouteFailurePolicy = policy
	}
	return params
}
