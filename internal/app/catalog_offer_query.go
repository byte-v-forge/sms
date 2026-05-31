package app

import (
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func normalizeOfferQuery(query core.RouteOfferQuery) core.RouteOfferQuery {
	query.ApplicationKey = strings.TrimSpace(query.ApplicationKey)
	query.CountryISO2 = strings.ToUpper(strings.TrimSpace(query.CountryISO2))
	query.CountryCallingCode = strings.TrimPrefix(strings.TrimSpace(query.CountryCallingCode), "+")
	query.ProviderKey = normalizeProviderKey(query.ProviderKey)
	return query
}

func routeOfferMatches(offer core.RouteOffer, query core.RouteOfferQuery) bool {
	if query.ProviderKey != "" && !strings.EqualFold(offer.ProviderKey, query.ProviderKey) {
		return false
	}
	if query.ApplicationKey != "" && !routeApplicationMatches(offer, query.ApplicationKey) {
		return false
	}
	if query.CountryISO2 != "" && offer.CountryISO2 != "" && !strings.EqualFold(offer.CountryISO2, query.CountryISO2) {
		return false
	}
	if query.CountryCallingCode != "" && offer.CountryCallingCode != "" && strings.TrimPrefix(offer.CountryCallingCode, "+") != query.CountryCallingCode {
		return false
	}
	return true
}

func routeApplicationMatches(offer core.RouteOffer, queryApplicationKey string) bool {
	query := normalizeCatalogToken(queryApplicationKey)
	if query == "" {
		return true
	}
	for _, candidate := range []string{offer.ApplicationKey, offer.Route.ApplicationKey, offer.Route.UpstreamServiceKey, offer.ApplicationName} {
		normalized := normalizeCatalogToken(candidate)
		if normalized == "" {
			continue
		}
		if normalized == query || strings.Contains(normalized, query) {
			return true
		}
	}
	return false
}

func normalizeCatalogToken(value string) string {
	var out strings.Builder
	for _, r := range strings.ToLower(strings.TrimSpace(value)) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			out.WriteRune(r)
		}
	}
	return out.String()
}
