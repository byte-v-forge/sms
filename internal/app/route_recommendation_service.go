package app

import (
	"context"
)

func (s *CatalogService) RecommendRoutes(ctx context.Context, query RouteRecommendationQuery) ([]RouteRecommendation, error) {
	target := normalizeRecommendationTarget(query.Target)
	if err := validateRecommendationTarget(target); err != nil {
		return nil, err
	}
	minPrice, maxPrice, err := recommendationPriceRange(query.Policy)
	if err != nil {
		return nil, err
	}
	providerFilter := normalizedProviderFilter(query.ProviderKeys)
	offers, err := s.listRecommendationOffers(ctx, target, providerFilter)
	if err != nil {
		return nil, err
	}
	strategy := recommendationStrategy(query.Policy)
	failurePolicy := routeFailurePolicyFromRoutePolicy(query.Policy)
	disabledRoutes := s.disabledRouteKeys(ctx, offers, failurePolicy)
	candidates := routeCandidates(offers, providerFilter, minPrice, maxPrice, disabledRoutes, failurePolicy)
	candidates = routeCandidatesWithMinimumAvailability(candidates, routeMinimumAvailability(query.Policy))
	if len(offers) > 0 && len(candidates) == 0 {
		return nil, routeRecommendationUnavailableError(target)
	}
	scoreRouteCandidates(candidates, strategy)
	sortRouteCandidates(candidates, strategy)
	return routeRecommendations(candidates, recommendationLimit(query.Policy)), nil
}
