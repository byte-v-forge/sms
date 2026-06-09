package fivesim

import (
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func fiveSimCountryIDForQuery(query core.RouteOfferQuery, countries []Country) string {
	for _, country := range countries {
		if query.CountryISO2 != "" && !strings.EqualFold(query.CountryISO2, country.CountryISO2) {
			continue
		}
		if query.CountryCallingCode != "" && strings.TrimPrefix(country.CountryCallingCode, "+") != strings.TrimPrefix(query.CountryCallingCode, "+") {
			continue
		}
		return country.CountryID
	}
	return ""
}
