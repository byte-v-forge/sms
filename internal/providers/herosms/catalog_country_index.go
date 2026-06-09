package herosms

import "strings"

func heroSMSCountryByID(countries []countryMetadata, countryID string) countryMetadata {
	for _, country := range countries {
		if country.ID == strings.TrimSpace(countryID) {
			return country
		}
	}
	return countryMetadata{ID: strings.TrimSpace(countryID)}
}
