package smsbower

import "github.com/byte-v-forge/sms/internal/core"

func smsbowerRoute(offer PriceOffer, iso2 string, callingCode string) core.Route {
	return core.Route{
		ProviderKey:        ProviderKey,
		ApplicationKey:     offer.UpstreamServiceKey,
		UpstreamServiceKey: offer.UpstreamServiceKey,
		CountryISO2:        iso2,
		CountryCallingCode: callingCode,
		ProviderCountryID:  offer.CountryID,
		UpstreamProviderID: offer.ProviderID,
	}
}
