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
