package geox

import (
	"strings"
	"sync"

	"github.com/biter777/countries"
)

var countryIndex = struct {
	sync.Once
	byAlpha map[string]countries.CountryCode
}{}

func countryByAlpha(value string) countries.CountryCode {
	value = strings.ToUpper(strings.TrimSpace(value))
	if value == "" {
		return countries.Unknown
	}
	countryIndex.Do(func() {
		countryIndex.byAlpha = make(map[string]countries.CountryCode, countries.Total()*2)
		for _, country := range countries.All() {
			if alpha2 := country.Alpha2(); alpha2 != "" {
				countryIndex.byAlpha[strings.ToUpper(alpha2)] = country
			}
			if alpha3 := country.Alpha3(); alpha3 != "" {
				countryIndex.byAlpha[strings.ToUpper(alpha3)] = country
			}
		}
	})
	return countryIndex.byAlpha[value]
}
