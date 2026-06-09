package app

import "github.com/byte-v-forge/sms/internal/core"

func routeCandidates(offers []core.RouteOffer, providerFilter map[string]struct{}, minPrice core.Money, maxPrice core.Money, disabledRoutes map[string]struct{}, failurePolicy core.RouteFailurePolicy) []routeCandidate {
	candidates := make([]routeCandidate, 0, len(offers))
	for _, offer := range offers {
		candidate, ok := routeCandidateForOffer(offer, providerFilter, minPrice, maxPrice, disabledRoutes, failurePolicy)
		if !ok {
			continue
		}
		candidates = append(candidates, candidate)
	}
	return candidates
}

func routeCandidateForOffer(offer core.RouteOffer, providerFilter map[string]struct{}, minPrice core.Money, maxPrice core.Money, disabledRoutes map[string]struct{}, failurePolicy core.RouteFailurePolicy) (routeCandidate, bool) {
	offer.Route = routeWithFailurePolicy(offer.Route, failurePolicy)
	if offer.AvailableCount <= 0 {
		return routeCandidate{}, false
	}
	if routeTemporarilyDisabled(offer.Route, disabledRoutes) {
		return routeCandidate{}, false
	}
	if !providerIncluded(offer.ProviderKey, providerFilter) {
		return routeCandidate{}, false
	}
	if !withinPriceRange(offer.Price, minPrice, maxPrice) {
		return routeCandidate{}, false
	}
	price, hasPrice := parseDecimalAmount(offer.Price.AmountDecimal)
	return routeCandidate{offer: offer, price: price, hasPrice: hasPrice}, true
}
