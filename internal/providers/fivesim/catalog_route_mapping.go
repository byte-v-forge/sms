package fivesim

import (
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func fiveSimRouteOffers(priceOffers []PriceOffer, countryByID map[string]Country) []core.RouteOffer {
	out := make([]core.RouteOffer, 0, len(priceOffers))
	now := time.Now().UTC()
	for _, offer := range priceOffers {
		out = append(out, fiveSimRouteOffer(offer, countryByID[offer.CountryID], now))
	}
	return out
}

func fiveSimRouteOffer(offer PriceOffer, country Country, observedAt time.Time) core.RouteOffer {
	return core.RouteOffer{
		ProviderKey:          ProviderKey,
		UpstreamProviderID:   offer.Operator,
		UpstreamProviderName: offer.Operator,
		ApplicationKey:       offer.UpstreamServiceKey,
		ApplicationName:      offer.UpstreamServiceKey,
		CountryISO2:          country.CountryISO2,
		CountryName:          country.Name,
		CountryCallingCode:   country.CountryCallingCode,
		Price:                offer.Price,
		AvailableCount:       offer.AvailableCount,
		ObservedAt:           observedAt,
		Route:                fiveSimRoute(offer, country),
	}
}

func fiveSimRoute(offer PriceOffer, country Country) core.Route {
	return core.Route{
		ProviderKey:        ProviderKey,
		ApplicationKey:     offer.UpstreamServiceKey,
		UpstreamServiceKey: offer.UpstreamServiceKey,
		CountryISO2:        country.CountryISO2,
		CountryCallingCode: country.CountryCallingCode,
		ProviderCountryID:  offer.CountryID,
		UpstreamProviderID: offer.Operator,
	}
}
