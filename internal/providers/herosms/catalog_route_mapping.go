package herosms

import (
	"time"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/stringx"
)

func heroSMSRouteOffersFromPrices(query core.RouteOfferQuery, priceOffers []PriceOffer, countries []countryMetadata, queryCountry countryMetadata, observedAt time.Time) []core.RouteOffer {
	out := make([]core.RouteOffer, 0, len(priceOffers))
	for _, offer := range priceOffers {
		if !heroSMSApplicationMatches(offer.UpstreamServiceKey, query.ApplicationKey) {
			continue
		}
		metadata := queryCountry
		if metadata.ID == "" {
			metadata = heroSMSCountryByID(countries, offer.CountryID)
		}
		out = append(out, heroSMSRouteOffer(query, offer, metadata, observedAt))
	}
	return out
}

func heroSMSRouteOffer(query core.RouteOfferQuery, offer PriceOffer, metadata countryMetadata, observedAt time.Time) core.RouteOffer {
	applicationKey := stringx.FirstNonEmpty(query.ApplicationKey, heroSMSPublicApplicationKey(offer.UpstreamServiceKey), offer.UpstreamServiceKey)
	route := core.Route{
		ProviderKey:        ProviderKey,
		ApplicationKey:     applicationKey,
		UpstreamServiceKey: offer.UpstreamServiceKey,
		CountryISO2:        metadata.ISO2,
		CountryCallingCode: metadata.CallingCode,
		ProviderCountryID:  offer.CountryID,
		MaxPrice:           offer.Price,
	}
	return core.RouteOffer{
		ProviderKey:          ProviderKey,
		UpstreamProviderName: "any",
		ApplicationKey:       applicationKey,
		ApplicationName:      heroSMSApplicationName(offer.UpstreamServiceKey),
		CountryISO2:          metadata.ISO2,
		CountryName:          stringx.FirstNonEmpty(metadata.Name, offer.CountryID),
		CountryCallingCode:   metadata.CallingCode,
		Price:                offer.Price,
		AvailableCount:       offer.AvailableCount,
		ObservedAt:           observedAt,
		Route:                route,
	}
}
