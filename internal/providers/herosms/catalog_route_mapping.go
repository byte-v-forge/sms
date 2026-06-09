package herosms

import (
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func heroSMSRouteOffersFromPrices(query core.RouteOfferQuery, priceOffers []PriceOffer, countries []countryMetadata, queryCountry countryMetadata, serviceNames map[string]string, observedAt time.Time) []core.RouteOffer {
	out := make([]core.RouteOffer, 0, len(priceOffers))
	for _, offer := range priceOffers {
		if !heroSMSApplicationMatches(offer.UpstreamServiceKey, query.ApplicationKey) {
			continue
		}
		out = append(out, heroSMSRouteOffer(query, offer, heroSMSOfferCountry(countries, queryCountry, offer), serviceNames, observedAt))
	}
	return out
}
