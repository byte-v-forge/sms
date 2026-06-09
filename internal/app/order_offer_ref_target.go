package app

import (
	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

func publicOfferTarget(route core.Route) *smsv1.SmsTarget {
	return &smsv1.SmsTarget{
		ApplicationKey:     routeText(route.ApplicationKey),
		CountryIso2:        routeCountryISO2(route.CountryISO2),
		CountryCallingCode: routeCallingCode(route.CountryCallingCode),
	}
}

func publicOfferRouteRef(route core.Route) *smsv1.SmsOfferRouteRef {
	return &smsv1.SmsOfferRouteRef{
		UpstreamServiceKey: routeText(route.UpstreamServiceKey),
		ProviderCountryId:  routeText(route.ProviderCountryID),
		UpstreamProviderId: routeText(route.UpstreamProviderID),
	}
}
