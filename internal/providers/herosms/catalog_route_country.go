package herosms

func heroSMSOfferCountry(countries []countryMetadata, queryCountry countryMetadata, offer PriceOffer) countryMetadata {
	if queryCountry.ID != "" {
		return queryCountry
	}
	return heroSMSCountryByID(countries, offer.CountryID)
}
