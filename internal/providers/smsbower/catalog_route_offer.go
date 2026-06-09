package smsbower

import (
	"time"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/stringx"
)

func smsbowerRouteOffer(offer PriceOffer, appNames map[string]string, country Country, observedAt time.Time) core.RouteOffer {
	iso2, callingCode := smsbowerCountryCodes(country.Name)
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
		Route:                smsbowerRoute(offer, iso2, callingCode),
	}
}
