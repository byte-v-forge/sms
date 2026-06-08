package geox

import (
	"strings"

	"github.com/biter777/countries"
)

func NormalizeCountryAlpha2(value string) string {
	country := countryByAlpha(value)
	if !country.IsValid() {
		return ""
	}
	return country.Alpha2()
}

func CountryAlpha2ByName(value string) string {
	return alpha2ByCountryName(value, false)
}

func CountryAlpha2InText(value string) string {
	return alpha2ByCountryName(value, true)
}

func CountryCallingCode(value string) string {
	country := countryByAlpha(value)
	if !country.IsValid() {
		country = countries.ByName(strings.TrimSpace(value))
	}
	if !country.IsValid() {
		return ""
	}
	for _, code := range country.CallCodes() {
		if code.IsValid() {
			return strings.TrimPrefix(code.String(), "+")
		}
	}
	return ""
}
