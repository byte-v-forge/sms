package app

import "strings"

func routeCandidateKey(candidate routeCandidate) string {
	offer := candidate.offer
	return strings.Join([]string{
		normalizeProviderKey(offer.ProviderKey),
		routeText(offer.ApplicationKey),
		routeCountryISO2(offer.CountryISO2),
		routeCallingCode(offer.CountryCallingCode),
		routeText(offer.UpstreamProviderID),
		routeText(offer.UpstreamProviderName),
		routeText(offer.Price.CurrencyCode),
		routeText(offer.Price.AmountDecimal),
	}, "\x00")
}
