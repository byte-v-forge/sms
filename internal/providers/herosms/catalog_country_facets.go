package herosms

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (c *Client) ListCatalogCountries(ctx context.Context, applicationKey string) ([]core.CatalogCountry, error) {
	countries, err := c.listCatalogCountryMetadata(ctx, applicationKey)
	if err != nil {
		return nil, err
	}
	out := make([]core.CatalogCountry, 0, len(countries))
	for _, country := range countries {
		out = append(out, core.CatalogCountry{
			CountryISO2:        country.ISO2,
			Name:               country.Name,
			CountryCallingCode: country.CallingCode,
		})
	}
	return out, nil
}

func (c *Client) listCatalogCountryMetadata(ctx context.Context, applicationKey string) ([]countryMetadata, error) {
	if normalizeHeroSMSServiceKey(applicationKey) == "" {
		return c.ListCountries(ctx)
	}
	services, err := c.ListServices(ctx)
	if err != nil {
		return nil, err
	}
	service := heroSMSServiceForQuery(applicationKey, services)
	if service == "" {
		return nil, nil
	}
	return c.ListServiceCountries(ctx, service)
}
