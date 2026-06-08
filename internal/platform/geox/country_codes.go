package geox

func CountryCodesInText(value string) []string {
	seen := map[string]struct{}{}
	out := []string{}
	addCountry := func(value string) {
		country := NormalizeCountryAlpha2(value)
		if country == "" {
			return
		}
		if _, exists := seen[country]; exists {
			return
		}
		seen[country] = struct{}{}
		out = append(out, country)
	}
	for _, value := range countryCodesFromEmoji(value) {
		addCountry(value)
	}
	for _, value := range alphaTokens(value) {
		addCountry(value)
	}
	return out
}
