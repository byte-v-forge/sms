package herosms

import (
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func heroSMSCountryForQuery(query core.RouteOfferQuery) countryMetadata {
	for _, country := range heroSMSCountries {
		if query.CountryISO2 != "" && !strings.EqualFold(query.CountryISO2, country.ISO2) {
			continue
		}
		if query.CountryCallingCode != "" && strings.TrimPrefix(query.CountryCallingCode, "+") != country.CallingCode {
			continue
		}
		return country
	}
	return countryMetadata{}
}

func heroSMSCountryByID(countryID string) countryMetadata {
	for _, country := range heroSMSCountries {
		if country.ID == strings.TrimSpace(countryID) {
			return country
		}
	}
	return countryMetadata{ID: strings.TrimSpace(countryID)}
}
