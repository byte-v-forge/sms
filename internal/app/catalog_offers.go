package app

import (
	"context"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
	"golang.org/x/sync/errgroup"
)

const catalogProviderConcurrency = 4

type RouteOfferList struct {
	Offers         []core.RouteOffer
	ProviderErrors []ProviderLookupError
}

type ProviderLookupError struct {
	ProviderKey         string
	ProviderDisplayName string
	Err                 error
}

func (s *CatalogService) ListPriceOffers(ctx context.Context, query core.RouteOfferQuery) ([]core.RouteOffer, error) {
	result, err := s.ListPriceOffersDetailed(ctx, query)
	if err != nil {
		return nil, err
	}
	return result.Offers, nil
}

func (s *CatalogService) ListPriceOffersDetailed(ctx context.Context, query core.RouteOfferQuery) (RouteOfferList, error) {
	query = normalizeOfferQuery(query)
	configs, err := s.configs.ListProviderConfigs(ctx, false, singleProviderKey(query.ProviderKeys))
	if err != nil {
		return RouteOfferList{}, err
	}
	configs = filteredCatalogProviderConfigs(configs, normalizedProviderFilter(query.ProviderKeys))
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
	var out []core.RouteOffer
	var providerErrors []ProviderLookupError
	var lastErr error
	for _, result := range results {
		if result.err != nil {
			lastErr = result.err
			providerErrors = append(providerErrors, ProviderLookupError{
				ProviderKey:         result.providerKey,
				ProviderDisplayName: result.providerDisplayName,
				Err:                 result.err,
			})
			continue
		}
		out = append(out, result.offers...)
	}
	if len(out) == 0 && lastErr != nil {
		return RouteOfferList{ProviderErrors: providerErrors}, lastErr
	}
	return RouteOfferList{Offers: out, ProviderErrors: providerErrors}, nil
}

type catalogProviderOffersResult struct {
	providerKey         string
	providerDisplayName string
	offers              []core.RouteOffer
	err                 error
}

func (s *CatalogService) listProviderPriceOffers(ctx context.Context, config *smsinternalv1.SmsProviderConfig, query core.RouteOfferQuery) ([]core.RouteOffer, error) {
	if !config.GetEnabled() {
		return nil, nil
	}
	provider, err := providerFromConfig(s.providers, config, s.timeout, s.defaultHTTPProxy)
	if err != nil {
		return nil, err
	}
	offerProvider, ok := provider.(routeOfferProvider)
	if !ok {
		return nil, nil
	}
	offers, err := offerProvider.ListRouteOffers(ctx, query)
	if err != nil {
		return nil, err
	}
	out := make([]core.RouteOffer, 0, len(offers))
	for _, offer := range offers {
		offer = s.finalizeOffer(config, offer)
		if !routeOfferMatches(offer, query) {
			continue
		}
		out = append(out, offer)
	}
	return out, nil
}

func filteredCatalogProviderConfigs(configs []*smsinternalv1.SmsProviderConfig, filter map[string]struct{}) []*smsinternalv1.SmsProviderConfig {
	if len(filter) == 0 {
		return configs
	}
	out := make([]*smsinternalv1.SmsProviderConfig, 0, len(configs))
	for _, config := range configs {
		if providerIncluded(config.GetProviderKey(), filter) {
			out = append(out, config)
		}
	}
	return out
}

func singleProviderKey(providerKeys []string) string {
	if len(providerKeys) == 1 {
		return providerKeys[0]
	}
	return ""
}
