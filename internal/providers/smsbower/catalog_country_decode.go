package smsbower

func decodeCountries(result string) ([]Country, error) {
	var raw map[string]countryShape
	if err := decodeJSONObject(result, &raw); err != nil {
		return nil, err
	}
	return countriesFromShapeMap(raw), nil
}
