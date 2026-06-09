package smsbower

import (
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func smsbowerCountryIDForQuery(query core.RouteOfferQuery, countries []Country) string {
	for _, country := range countries {
		iso2, callingCode := smsbowerCountryCodes(country.Name)
		if query.CountryISO2 != "" && !strings.EqualFold(query.CountryISO2, iso2) {
			continue
		}
		if query.CountryCallingCode != "" && query.CountryCallingCode != strings.TrimPrefix(callingCode, "+") {
			continue
		}
		return country.CountryID
	}
	return ""
}
