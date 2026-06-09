package fivesim

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (c *Client) ListCatalogCountries(ctx context.Context, _ string) ([]core.CatalogCountry, error) {
	countries, err := c.ListCountries(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]core.CatalogCountry, 0, len(countries))
	for _, country := range countries {
		out = append(out, core.CatalogCountry{
			CountryISO2:        country.CountryISO2,
			Name:               country.Name,
			CountryCallingCode: country.CountryCallingCode,
		})
	}
	return out, nil
}
