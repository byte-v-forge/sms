package herosms

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/geox"
	"github.com/byte-v-forge/sms/internal/platform/jsonx"
	"github.com/byte-v-forge/sms/internal/platform/stringx"
)

func (c *Client) ListCountries(ctx context.Context) ([]countryMetadata, error) {
	result, err := c.api.Do(ctx, "getCountries", nil)
	if err != nil {
		return nil, err
	}
	var raw map[string]struct {
		ID   json.RawMessage `json:"id"`
		Name string          `json:"name"`
		Rus  string          `json:"rus"`
		Eng  string          `json:"eng"`
		Chn  string          `json:"chn"`
	}
	if err := decodeHeroSMSJSONObject(result, &raw); err != nil {
		return nil, err
	}
	countries := make([]countryMetadata, 0, len(raw))
	for key, item := range raw {
		id := jsonx.FirstScalar(item.ID)
		if id == "" || id == "0" {
			id = key
		}
		name := stringx.FirstNonEmpty(item.Eng, item.Name, item.Chn, item.Rus, key)
		iso2 := geox.CountryAlpha2InText(name)
		countries = append(countries, countryMetadata{ID: id, Name: name, ISO2: iso2, CallingCode: geox.CountryCallingCode(iso2)})
	}
	return countries, nil
}

func heroSMSCountryForQuery(query core.RouteOfferQuery, countries []countryMetadata) countryMetadata {
	for _, country := range countries {
		if query.CountryISO2 != "" && !strings.EqualFold(query.CountryISO2, country.ISO2) {
			continue
		}
		if query.CountryCallingCode != "" && strings.TrimPrefix(query.CountryCallingCode, "+") != country.CallingCode {
			continue
		}
		return country
	}
	return countryMetadata{}
}

func heroSMSCountryByID(countries []countryMetadata, countryID string) countryMetadata {
	for _, country := range countries {
		if country.ID == strings.TrimSpace(countryID) {
			return country
		}
	}
	return countryMetadata{ID: strings.TrimSpace(countryID)}
}
