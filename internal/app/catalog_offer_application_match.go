package app

import "github.com/byte-v-forge/sms/internal/core"

func routeApplicationMatches(offer core.RouteOffer, queryApplicationKey string) bool {
	query := normalizeCatalogToken(queryApplicationKey)
	if query == "" {
		return true
	}
	for _, candidate := range routeApplicationCandidates(offer) {
		if catalogTokenMatches(candidate, query) {
			return true
		}
	}
	return false
}

func routeApplicationCandidates(offer core.RouteOffer) []string {
	return []string{offer.ApplicationKey, offer.Route.ApplicationKey, offer.Route.UpstreamServiceKey, offer.ApplicationName}
}
