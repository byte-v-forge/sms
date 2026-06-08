package app

import (
	"math/big"

	"github.com/byte-v-forge/sms/internal/core"
)

type routeCandidate struct {
	offer    core.RouteOffer
	price    *big.Rat
	hasPrice bool
	score    int32
}

func routeCandidates(offers []core.RouteOffer, providerFilter map[string]struct{}, minPrice core.Money, maxPrice core.Money, disabledRoutes map[string]struct{}, failurePolicy core.RouteFailurePolicy) []routeCandidate {
	candidates := make([]routeCandidate, 0, len(offers))
	for _, offer := range offers {
		offer.Route = routeWithFailurePolicy(offer.Route, failurePolicy)
		if offer.AvailableCount <= 0 {
			continue
		}
		if routeTemporarilyDisabled(offer.Route, disabledRoutes) {
			continue
		}
		if !providerIncluded(offer.ProviderKey, providerFilter) {
			continue
		}
		if !withinPriceRange(offer.Price, minPrice, maxPrice) {
			continue
		}
		price, hasPrice := parseDecimalAmount(offer.Price.AmountDecimal)
		candidates = append(candidates, routeCandidate{offer: offer, price: price, hasPrice: hasPrice})
	}
	return candidates
}

func providerIncluded(providerKey string, filter map[string]struct{}) bool {
	if len(filter) == 0 {
		return true
	}
	_, ok := filter[normalizeProviderKey(providerKey)]
	return ok
}

func routeRecommendations(candidates []routeCandidate, limit int) []RouteRecommendation {
	if limit > len(candidates) {
		limit = len(candidates)
	}
	recommendations := make([]RouteRecommendation, 0, limit)
	for index := 0; index < limit; index++ {
		recommendations = append(recommendations, RouteRecommendation{
			Offer: candidates[index].offer,
			Score: candidates[index].score,
		})
	}
	return recommendations
}
