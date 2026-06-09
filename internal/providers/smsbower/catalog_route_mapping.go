package smsbower

import (
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func smsbowerRouteOffers(priceOffers []PriceOffer, catalog routeCatalogInputs) []core.RouteOffer {
	appNames := smsbowerApplicationNames(catalog.applications)
	countryByID := smsbowerCountriesByID(catalog.countries)
	out := make([]core.RouteOffer, 0, len(priceOffers))
	now := time.Now().UTC()
	for _, offer := range priceOffers {
		out = append(out, smsbowerRouteOffer(offer, appNames, countryByID[offer.CountryID], now))
	}
	return out
}
