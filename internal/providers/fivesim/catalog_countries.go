package fivesim

import (
	"context"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func (c *Client) ListCountries(ctx context.Context) ([]Country, error) {
	var raw map[string]struct {
		ISO    map[string]int `json:"iso"`
		Prefix map[string]int `json:"prefix"`
		Name   string         `json:"text_en"`
	}
	if err := c.getJSON(ctx, "/v1/guest/countries", nil, false, &raw); err != nil {
		return nil, err
	}
	countries := make([]Country, 0, len(raw))
	for countryID, item := range raw {
		countries = append(countries, Country{
			CountryID:          countryID,
			Name:               item.Name,
			CountryISO2:        firstMapKey(item.ISO),
			CountryCallingCode: trimPlus(firstMapKey(item.Prefix)),
		})
	}
	return countries, nil
}
func fiveSimCountryIDForQuery(query core.RouteOfferQuery, countries []Country) string {
	for _, country := range countries {
		if query.CountryISO2 != "" && !strings.EqualFold(query.CountryISO2, country.CountryISO2) {
			continue
		}
		if query.CountryCallingCode != "" && strings.TrimPrefix(country.CountryCallingCode, "+") != strings.TrimPrefix(query.CountryCallingCode, "+") {
			continue
		}
		return country.CountryID
	}
	return ""
}

func firstMapKey(values map[string]int) string {
	for key := range values {
		return key
	}
	return ""
}

func trimPlus(value string) string {
	if len(value) > 0 && value[0] == '+' {
		return value[1:]
	}
	return value
}
