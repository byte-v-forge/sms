package app

import "github.com/byte-v-forge/sms/internal/core"

func normalizeOfferQuery(query core.RouteOfferQuery) core.RouteOfferQuery {
	query.ApplicationKey = routeText(query.ApplicationKey)
	query.CountryISO2 = routeCountryISO2(query.CountryISO2)
	query.CountryCallingCode = routeCallingCode(query.CountryCallingCode)
	query.ProviderKeys = sortedProviderFilterKeys(normalizedProviderFilter(query.ProviderKeys))
	return query
}
