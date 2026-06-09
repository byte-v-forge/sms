package herosms

func decodeHeroSMSCountries(result string) ([]countryMetadata, error) {
	var raw map[string]heroSMSCountryShape
	if err := decodeHeroSMSJSONObject(result, &raw); err != nil {
		return nil, err
	}
	return heroSMSCountriesFromShapeMap(raw), nil
}
