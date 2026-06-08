package app

import (
	"context"
	"math/big"
	"strings"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

type routeCandidate struct {
	offer    core.RouteOffer
	price    *big.Rat
	hasPrice bool
	score    int32
}

func (s *CatalogService) disabledRouteKeys(ctx context.Context, offers []core.RouteOffer, failurePolicy core.RouteFailurePolicy) map[string]struct{} {
	if len(offers) == 0 || s.routeHealth == nil {
		return map[string]struct{}{}
	}
	routes := make([]core.Route, 0, len(offers))
	for _, offer := range offers {
		routes = append(routes, routeWithFailurePolicy(offer.Route, failurePolicy))
	}
	disabled, err := s.routeHealth.DisabledRouteKeys(ctx, routes)
	if err != nil {
		return map[string]struct{}{}
	}
	return disabled
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

func routeWithFailurePolicy(route core.Route, failurePolicy core.RouteFailurePolicy) core.Route {
	route.FailurePolicy = failurePolicy
	return route
}

func routeTemporarilyDisabled(route core.Route, disabledRoutes map[string]struct{}) bool {
	if len(disabledRoutes) == 0 {
		return false
	}
	_, disabled := disabledRoutes[routeHealthKey(route)]
	return disabled
}

func routeMinimumAvailability(policy *smsv1.SmsRoutePolicy) int {
	if policy.GetMinAvailableCount() > 0 {
		return int(policy.GetMinAvailableCount())
	}
	return 1
}

func routeCandidatesWithMinimumAvailability(candidates []routeCandidate, minimum int) []routeCandidate {
	if minimum <= 1 {
		return candidates
	}
	out := make([]routeCandidate, 0, len(candidates))
	for _, candidate := range candidates {
		if candidate.offer.AvailableCount >= minimum {
			out = append(out, candidate)
		}
	}
	return out
}

func providerIncluded(providerKey string, filter map[string]struct{}) bool {
	if len(filter) == 0 {
		return true
	}
	_, ok := filter[normalizeProviderKey(providerKey)]
	return ok
}

func withinPriceRange(price core.Money, minPrice core.Money, maxPrice core.Money) bool {
	offerAmount, ok := parseDecimalAmount(price.AmountDecimal)
	if !ok {
		return !moneyIsSet(minPrice) && !moneyIsSet(maxPrice)
	}
	if !withinPriceBound(price, offerAmount, minPrice, 1) {
		return false
	}
	return withinPriceBound(price, offerAmount, maxPrice, -1)
}

func withinPriceBound(price core.Money, offerAmount *big.Rat, bound core.Money, direction int) bool {
	if strings.TrimSpace(bound.AmountDecimal) == "" {
		return true
	}
	if bound.CurrencyCode != "" && price.CurrencyCode != "" && !strings.EqualFold(bound.CurrencyCode, price.CurrencyCode) {
		return false
	}
	boundAmount, ok := parseDecimalAmount(bound.AmountDecimal)
	if !ok {
		return false
	}
	comparison := offerAmount.Cmp(boundAmount)
	if direction > 0 {
		return comparison >= 0
	}
	return comparison <= 0
}

func parseDecimalAmount(value string) (*big.Rat, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, false
	}
	amount, ok := new(big.Rat).SetString(value)
	if !ok {
		return nil, false
	}
	return amount, true
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
