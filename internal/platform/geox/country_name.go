package geox

import (
	"github.com/biter777/countries"
)

func alpha2ByCountryName(value string, inText bool) string {
	normalized := normalizeCountryNameText(value)
	if normalized == "" {
		return ""
	}
	country := countries.ByName(normalized)
	if country.IsValid() {
		return country.Alpha2()
	}
	for _, countryName := range indexedCountryNames() {
		if countryName.name == "" {
			continue
		}
		if normalized == countryName.name || inText && countryNameOccursInText(normalized, countryName.name) {
			return countryName.alpha2
		}
	}
	return ""
}
