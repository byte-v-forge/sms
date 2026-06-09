package fivesim

func fiveSimCountriesByID(countries []Country) map[string]Country {
	countryByID := map[string]Country{}
	for _, country := range countries {
		countryByID[country.CountryID] = country
	}
	return countryByID
}
