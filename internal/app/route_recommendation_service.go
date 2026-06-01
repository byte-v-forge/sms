package app

import (
	"context"
	"fmt"

	smsv1 "github.com/byte-v-forge/common-lib/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

const (
	defaultRouteRecommendationLimit = 20
	maxRouteRecommendationLimit     = 100
)

type RouteRecommendationQuery struct {
	Target       core.Target
	Policy       *smsv1.SmsRoutePolicy
	ProviderKeys []string
}

type RouteRecommendation struct {
	Offer core.RouteOffer
	Score int32
}

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

func routeRecommendationUnavailableError(target core.Target) error {
	return core.NewError(
		core.CodeSupplyUnavailable,
		fmt.Sprintf("sms route currently unavailable for %s/%s/%s", target.ApplicationKey, target.CountryISO2, target.CountryCallingCode),
		true,
	)
}

func (s *CatalogService) listRecommendationOffers(ctx context.Context, target core.Target, providerFilter map[string]struct{}) ([]core.RouteOffer, error) {
	query := core.RouteOfferQuery{
		ApplicationKey:     target.ApplicationKey,
		CountryISO2:        target.CountryISO2,
		CountryCallingCode: target.CountryCallingCode,
	}
	if len(providerFilter) == 0 {
		return s.ListPriceOffers(ctx, query)
	}
	var out []core.RouteOffer
	var lastErr error
	for _, providerKey := range sortedProviderFilterKeys(providerFilter) {
		query.ProviderKey = providerKey
		offers, err := s.ListPriceOffers(ctx, query)
		if err != nil {
			lastErr = err
			continue
		}
		out = append(out, offers...)
	}
	if len(out) == 0 && lastErr != nil {
		return nil, lastErr
	}
	return out, nil
}
