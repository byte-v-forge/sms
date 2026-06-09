package smsbower

import (
	"encoding/json"

	"github.com/byte-v-forge/sms/internal/platform/jsonx"
	"github.com/byte-v-forge/sms/internal/platform/stringx"
)

type countryShape struct {
	ID  json.RawMessage `json:"id"`
	Rus string          `json:"rus"`
	Eng string          `json:"eng"`
	Chn string          `json:"chn"`
}

func countriesFromShapeMap(raw map[string]countryShape) []Country {
	countries := make([]Country, 0, len(raw))
	for key, item := range raw {
		countries = append(countries, countryFromShape(key, item))
	}
	return countries
}

func countryFromShape(key string, item countryShape) Country {
	id := jsonx.Scalar(item.ID)
	if id == "" || id == "0" {
		id = key
	}
	name := stringx.FirstNonEmpty(item.Eng, item.Chn, item.Rus, key)
	return Country{CountryID: id, Name: name}
}
