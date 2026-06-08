package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *CatalogService) listRecommendationOffers(ctx context.Context, target core.Target, providerFilter map[string]struct{}) ([]core.RouteOffer, error) {
	query := core.RouteOfferQuery{
		ApplicationKey:     target.ApplicationKey,
		CountryISO2:        target.CountryISO2,
		CountryCallingCode: target.CountryCallingCode,
	}
	query.ProviderKeys = sortedProviderFilterKeys(providerFilter)
	return s.ListPriceOffers(ctx, query)
}
