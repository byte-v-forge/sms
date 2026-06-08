package app

import (
	"context"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

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
