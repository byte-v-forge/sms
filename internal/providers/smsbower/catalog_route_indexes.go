package smsbower

import "github.com/byte-v-forge/sms/internal/platform/stringx"

func smsbowerApplicationNames(applications []ApplicationOffer) map[string]string {
	appNames := map[string]string{}
	for _, app := range applications {
		appNames[app.UpstreamServiceKey] = stringx.FirstNonEmpty(app.DisplayName, app.UpstreamServiceKey)
	}
	return appNames
}

func smsbowerCountriesByID(countries []Country) map[string]Country {
	countryByID := map[string]Country{}
	for _, country := range countries {
		countryByID[country.CountryID] = country
	}
	return countryByID
}
