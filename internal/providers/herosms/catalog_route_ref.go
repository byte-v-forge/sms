package herosms

import "github.com/byte-v-forge/sms/internal/core"

func heroSMSRoute(applicationKey string, offer PriceOffer, metadata countryMetadata) core.Route {
	return core.Route{
		ProviderKey:        ProviderKey,
		ApplicationKey:     applicationKey,
		UpstreamServiceKey: offer.UpstreamServiceKey,
		CountryISO2:        metadata.ISO2,
		CountryCallingCode: metadata.CallingCode,
		ProviderCountryID:  offer.CountryID,
		UpstreamProviderID: heroSMSOperator(offer.Operator),
		MaxPrice:           offer.Price,
	}
}
