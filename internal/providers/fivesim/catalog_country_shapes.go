package fivesim

type countryShape struct {
	ISO    map[string]int `json:"iso"`
	Prefix map[string]int `json:"prefix"`
	Name   string         `json:"text_en"`
}

func fiveSimCountriesFromShapes(raw map[string]countryShape) []Country {
	countries := make([]Country, 0, len(raw))
	for countryID, item := range raw {
		countries = append(countries, fiveSimCountryFromShape(countryID, item))
	}
	return countries
}

func fiveSimCountryFromShape(countryID string, item countryShape) Country {
	return Country{
		CountryID:          countryID,
		Name:               item.Name,
		CountryISO2:        firstMapKey(item.ISO),
		CountryCallingCode: trimPlus(firstMapKey(item.Prefix)),
	}
}
