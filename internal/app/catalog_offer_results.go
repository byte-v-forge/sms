package app

import "github.com/byte-v-forge/sms/internal/core"

func collectCatalogProviderOffers(results []catalogProviderOffersResult) (RouteOfferList, error) {
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
