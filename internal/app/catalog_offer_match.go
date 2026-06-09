package app

import (
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func routeOfferMatches(offer core.RouteOffer, query core.RouteOfferQuery) bool {
	if !providerIncluded(offer.ProviderKey, normalizedProviderFilter(query.ProviderKeys)) {
		return false
	}
	if query.ApplicationKey != "" && !routeApplicationMatches(offer, query.ApplicationKey) {
		return false
	}
	if query.CountryISO2 != "" && offer.CountryISO2 != "" && !strings.EqualFold(offer.CountryISO2, query.CountryISO2) {
		return false
	}
	if query.CountryCallingCode != "" && offer.CountryCallingCode != "" && routeCallingCode(offer.CountryCallingCode) != query.CountryCallingCode {
		return false
	}
	return true
}
