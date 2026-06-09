package herosms

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/geox"
	"github.com/byte-v-forge/sms/internal/platform/jsonx"
)

type serviceCountriesResponse struct {
	Data json.RawMessage `json:"data"`
}

type serviceCountryShape struct {
	ID   json.RawMessage `json:"id"`
	Name string          `json:"name"`
}

func (c *Client) ListServiceCountries(ctx context.Context, serviceKey string) ([]countryMetadata, error) {
	path := "/left-menu/service/" + url.PathEscape(normalizeHeroSMSServiceKey(serviceKey)) + "/countries"
	var response serviceCountriesResponse
	if err := c.getOpenAPIJSON(ctx, path, nil, &response); err != nil {
		return nil, err
	}
	return serviceCountriesFromRaw(response.Data)
}

func serviceCountriesFromRaw(raw json.RawMessage) ([]countryMetadata, error) {
	if countries, ok := serviceCountriesFromList(raw); ok {
		return countries, nil
	}
	if countries, ok := serviceCountriesFromMap(raw); ok {
		return countries, nil
	}
	return nil, core.NewError(core.CodeUpstreamRejected, "bad hero sms service countries response", false)
}

func serviceCountriesFromList(raw json.RawMessage) ([]countryMetadata, bool) {
	var shapes []serviceCountryShape
	if err := json.Unmarshal(raw, &shapes); err != nil {
		return nil, false
	}
	return serviceCountriesFromShapes(shapes), true
}

func serviceCountriesFromMap(raw json.RawMessage) ([]countryMetadata, bool) {
	var shapes map[string]serviceCountryShape
	if err := json.Unmarshal(raw, &shapes); err != nil {
		return nil, false
	}
	out := make([]serviceCountryShape, 0, len(shapes))
	for key, shape := range shapes {
		if jsonx.FirstScalar(shape.ID) == "" {
			shape.ID, _ = json.Marshal(key)
		}
		out = append(out, shape)
	}
	return serviceCountriesFromShapes(out), true
}

func serviceCountriesFromShapes(shapes []serviceCountryShape) []countryMetadata {
	countries := make([]countryMetadata, 0, len(shapes))
	for _, shape := range shapes {
		iso2 := geox.CountryAlpha2InText(shape.Name)
		countries = append(countries, countryMetadata{
			ID:          jsonx.FirstScalar(shape.ID),
			Name:        shape.Name,
			ISO2:        iso2,
			CallingCode: geox.CountryCallingCode(iso2),
		})
	}
	return countries
}
