package herosms

import (
	"time"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/stringx"
)

func heroSMSRouteOffer(query core.RouteOfferQuery, offer PriceOffer, metadata countryMetadata, observedAt time.Time) core.RouteOffer {
	applicationKey := heroSMSOfferApplicationKey(query, offer)
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
		Route:                heroSMSRoute(applicationKey, offer, metadata),
	}
}

func heroSMSOfferApplicationKey(query core.RouteOfferQuery, offer PriceOffer) string {
	return stringx.FirstNonEmpty(query.ApplicationKey, heroSMSPublicApplicationKey(offer.UpstreamServiceKey), offer.UpstreamServiceKey)
}
