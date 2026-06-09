package geox

func CountryRegionCode(value string) string {
	country := countryByAlpha(value)
	if !country.IsValid() {
		return ""
	}
	return regionShortCode(country.Region())
}
