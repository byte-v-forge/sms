package app

import "github.com/byte-v-forge/sms/internal/core"

func finalizeOfferRouteDefaults(offer core.RouteOffer) core.RouteOffer {
	if offer.Route.ApplicationKey == "" {
		offer.Route.ApplicationKey = offer.ApplicationKey
	}
	if offer.Route.CountryISO2 == "" {
		offer.Route.CountryISO2 = offer.CountryISO2
	}
	if offer.Route.CountryCallingCode == "" {
		offer.Route.CountryCallingCode = offer.CountryCallingCode
	}
	if offer.ApplicationKey == "" {
		offer.ApplicationKey = offer.Route.ApplicationKey
	}
	if offer.CountryISO2 == "" {
		offer.CountryISO2 = offer.Route.CountryISO2
	}
	if offer.CountryCallingCode == "" {
		offer.CountryCallingCode = offer.Route.CountryCallingCode
	}
	return offer
}
