package smsbower

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/byte-v-forge/common-lib/geox"
	"github.com/byte-v-forge/common-lib/stringx"
	"github.com/byte-v-forge/sms/internal/core"
)

func (c *Client) ListCountries(ctx context.Context) ([]Country, error) {
	result, err := c.api.Do(ctx, "getCountries", nil)
	if err != nil {
		return nil, err
	}
	var raw map[string]struct {
		ID  json.RawMessage `json:"id"`
		Rus string          `json:"rus"`
		Eng string          `json:"eng"`
		Chn string          `json:"chn"`
	}
	if err := decodeJSONObject(result, &raw); err != nil {
		return nil, err
	}
	countries := make([]Country, 0, len(raw))
	for key, item := range raw {
		id := rawJSONScalar(item.ID)
		if id == "" || id == "0" {
			id = key
		}
		name := stringx.FirstNonEmpty(item.Eng, item.Chn, item.Rus, key)
		countries = append(countries, Country{CountryID: id, Name: name})
	}
	return countries, nil
}
func smsbowerCountryIDForQuery(query core.RouteOfferQuery, countries []Country) string {
	for _, country := range countries {
		iso2, callingCode := smsbowerCountryCodes(country.Name)
		if query.CountryISO2 != "" && !strings.EqualFold(query.CountryISO2, iso2) {
			continue
		}
		if query.CountryCallingCode != "" && query.CountryCallingCode != strings.TrimPrefix(callingCode, "+") {
			continue
		}
		return country.CountryID
	}
	return ""
}

func smsbowerCountryCodes(name string) (string, string) {
	normalized := strings.TrimSpace(name)
	if normalized == "" {
		return "", ""
	}
	iso2 := geox.CountryAlpha2InText(normalized)
	if iso2 == "" {
		return "", ""
	}
	return iso2, geox.CountryCallingCode(iso2)
}
