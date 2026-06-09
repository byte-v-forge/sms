package herosms

import (
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func heroSMSCountryForQuery(query core.RouteOfferQuery, countries []countryMetadata) countryMetadata {
	for _, country := range countries {
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
