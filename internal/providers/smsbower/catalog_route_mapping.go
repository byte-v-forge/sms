package smsbower

import (
	"time"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/stringx"
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

func smsbowerRouteOffer(offer PriceOffer, appNames map[string]string, country Country, observedAt time.Time) core.RouteOffer {
	iso2, callingCode := smsbowerCountryCodes(country.Name)
	route := core.Route{
		ProviderKey:        ProviderKey,
		ApplicationKey:     offer.UpstreamServiceKey,
		UpstreamServiceKey: offer.UpstreamServiceKey,
		CountryISO2:        iso2,
		CountryCallingCode: callingCode,
		ProviderCountryID:  offer.CountryID,
		UpstreamProviderID: offer.ProviderID,
	}
	return core.RouteOffer{
		ProviderKey:          ProviderKey,
		UpstreamProviderID:   offer.ProviderID,
		UpstreamProviderName: offer.ProviderID,
		ApplicationKey:       offer.UpstreamServiceKey,
		ApplicationName:      stringx.FirstNonEmpty(appNames[offer.UpstreamServiceKey], offer.UpstreamServiceKey),
		CountryISO2:          iso2,
		CountryName:          stringx.FirstNonEmpty(country.Name, offer.CountryID),
		CountryCallingCode:   callingCode,
		Price:                offer.Price,
		AvailableCount:       offer.AvailableCount,
		ObservedAt:           observedAt,
		Route:                route,
	}
}
