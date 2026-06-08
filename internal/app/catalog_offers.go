package app

import (
	"context"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
	"golang.org/x/sync/errgroup"
)

const catalogProviderConcurrency = 4

func (s *CatalogService) ListPriceOffers(ctx context.Context, query core.RouteOfferQuery) ([]core.RouteOffer, error) {
	query = normalizeOfferQuery(query)
	configs, err := s.configs.ListProviderConfigs(ctx, false, singleProviderKey(query.ProviderKeys))
	if err != nil {
		return nil, err
	}
	configs = filteredCatalogProviderConfigs(configs, normalizedProviderFilter(query.ProviderKeys))
	results := make([]catalogProviderOffersResult, len(configs))
	group, groupCtx := errgroup.WithContext(ctx)
	group.SetLimit(catalogProviderConcurrency)
	for index, config := range configs {
		index := index
		config := config
		group.Go(func() error {
			results[index].offers, results[index].err = s.listProviderPriceOffers(groupCtx, config, query)
			return nil
		})
	}
	if err := group.Wait(); err != nil {
		return nil, err
	}
	var out []core.RouteOffer
	var lastErr error
	for _, result := range results {
		if result.err != nil {
			lastErr = result.err
			continue
		}
		out = append(out, result.offers...)
	}
	if len(out) == 0 && lastErr != nil {
		return nil, lastErr
	}
	return out, nil
}

type catalogProviderOffersResult struct {
	offers []core.RouteOffer
	err    error
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
