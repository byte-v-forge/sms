package app

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (s *CatalogService) ListPriceOffers(ctx context.Context, query core.RouteOfferQuery) ([]core.RouteOffer, error) {
	query = normalizeOfferQuery(query)
	configs, err := s.configs.ListProviderConfigs(ctx, false, query.ProviderKey)
	if err != nil {
		return nil, err
	}
	var out []core.RouteOffer
	var lastErr error
	for _, config := range configs {
		if !config.GetEnabled() {
			continue
		}
		provider, err := providerFromConfig(s.providers, config, s.timeout, s.defaultHTTPProxy)
		if err != nil {
			lastErr = err
			continue
		}
		offerProvider, ok := provider.(routeOfferProvider)
		if !ok {
			continue
		}
		offers, err := offerProvider.ListRouteOffers(ctx, query)
		if err != nil {
			lastErr = err
			continue
		}
		for _, offer := range offers {
			offer = s.finalizeOffer(config, offer)
			if !routeOfferMatches(offer, query) {
				continue
			}
			out = append(out, offer)
		}
	}
	if len(out) == 0 && lastErr != nil {
		return nil, lastErr
	}
	return out, nil
}
