package herosms

import (
	"encoding/json"

	"github.com/byte-v-forge/sms/internal/platform/geox"
	"github.com/byte-v-forge/sms/internal/platform/jsonx"
	"github.com/byte-v-forge/sms/internal/platform/stringx"
)

type heroSMSCountryShape struct {
	ID   json.RawMessage `json:"id"`
	Name string          `json:"name"`
	Rus  string          `json:"rus"`
	Eng  string          `json:"eng"`
	Chn  string          `json:"chn"`
}

func heroSMSCountriesFromShapeMap(raw map[string]heroSMSCountryShape) []countryMetadata {
	countries := make([]countryMetadata, 0, len(raw))
	for key, item := range raw {
		countries = append(countries, heroSMSCountryFromShape(key, item))
	}
	return countries
}

func heroSMSCountryFromShape(key string, item heroSMSCountryShape) countryMetadata {
	id := jsonx.FirstScalar(item.ID)
	if id == "" || id == "0" {
		id = key
	}
	name := stringx.FirstNonEmpty(item.Eng, item.Name, item.Chn, item.Rus, key)
	iso2 := geox.CountryAlpha2InText(name)
	return countryMetadata{
		ID:          id,
		Name:        name,
		ISO2:        iso2,
		CallingCode: geox.CountryCallingCode(iso2),
	}
}
