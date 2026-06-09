package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
	"golang.org/x/sync/errgroup"
)

func (s *CatalogService) ListPriceOffers(ctx context.Context, query core.RouteOfferQuery) ([]core.RouteOffer, error) {
	result, err := s.ListPriceOffersDetailed(ctx, query)
	if err != nil {
		return nil, err
	}
	return result.Offers, nil
}

func (s *CatalogService) ListPriceOffersDetailed(ctx context.Context, query core.RouteOfferQuery) (RouteOfferList, error) {
	query = normalizeOfferQuery(query)
	configs, err := s.catalogProviderConfigs(ctx, query.ProviderKeys)
	if err != nil {
		return RouteOfferList{}, err
	}
	results := make([]catalogProviderOffersResult, len(configs))
	group, groupCtx := errgroup.WithContext(ctx)
	group.SetLimit(catalogProviderConcurrency)
	for index, config := range configs {
		index := index
		config := config
		results[index].providerKey = normalizeProviderKey(config.GetProviderKey())
		results[index].providerDisplayName = s.providers.DisplayName(config.GetProviderKey())
		group.Go(func() error {
			results[index].offers, results[index].err = s.listProviderPriceOffers(groupCtx, config, query)
			return nil
		})
	}
	if err := group.Wait(); err != nil {
		return RouteOfferList{}, err
	}
	return collectCatalogProviderOffers(results)
}
