package app

import (
	"sort"

	"github.com/byte-v-forge/sms/internal/core"
)

func collectCatalogProviderCountries(results []catalogProviderCountriesResult) ([]core.CatalogCountry, error) {
	items := map[string]core.CatalogCountry{}
	lastErr := countryProviderError(results)
	for _, result := range results {
		for _, country := range result.countries {
			key := catalogCountryKey(country)
			if key == "" {
				continue
			}
			items[key] = bestCatalogCountry(items[key], core.CatalogCountry{
				CountryISO2:        routeCountryISO2(country.CountryISO2),
				Name:               firstNonEmpty(routeText(country.Name), key),
				CountryCallingCode: routeCallingCode(country.CountryCallingCode),
			})
		}
	}
	countries := catalogCountryValues(items)
	if len(countries) == 0 && lastErr != nil {
		return nil, lastErr
	}
	return countries, nil
}

func countryProviderError(results []catalogProviderCountriesResult) error {
	var lastErr error
	for _, result := range results {
		if result.err == nil {
			continue
		}
		lastErr = result.err
	}
	return lastErr
}

func catalogCountryKey(country core.CatalogCountry) string {
	iso2 := routeCountryISO2(country.CountryISO2)
	if iso2 != "" {
		return iso2
	}
	return routeCallingCode(country.CountryCallingCode)
}

func bestCatalogCountry(left core.CatalogCountry, right core.CatalogCountry) core.CatalogCountry {
	if left.Name == "" || len(right.Name) > len(left.Name) {
		return right
	}
	return left
}

func catalogCountryValues(items map[string]core.CatalogCountry) []core.CatalogCountry {
	out := make([]core.CatalogCountry, 0, len(items))
	for _, item := range items {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
