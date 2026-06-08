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
			ApplicationKey:     routeText(route.ApplicationKey),
			CountryIso2:        routeCountryISO2(route.CountryISO2),
			CountryCallingCode: routeCallingCode(route.CountryCallingCode),
		},
		RouteRef: &smsv1.SmsOfferRouteRef{
			UpstreamServiceKey: routeText(route.UpstreamServiceKey),
			ProviderCountryId:  routeText(route.ProviderCountryID),
			UpstreamProviderId: routeText(route.UpstreamProviderID),
		},
	}
	if ref.GetProviderKey() == "" && ref.GetOfferId() == "" && targetIsZero(ref.GetTarget()) && offerRouteRefIsZero(ref.GetRouteRef()) {
		return nil
	}
	return ref
}
