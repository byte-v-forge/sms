package smsbower

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (c *Client) ListCatalogApplications(ctx context.Context) ([]core.CatalogApplication, error) {
	applications, err := c.ListApplications(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]core.CatalogApplication, 0, len(applications))
	for _, app := range applications {
		out = append(out, core.CatalogApplication{
			ApplicationKey: app.ApplicationKey,
			DisplayName:    app.DisplayName,
			Aliases:        []string{app.ApplicationKey, app.UpstreamServiceKey, app.DisplayName},
		})
	}
	return out, nil
}

func (c *Client) ListCatalogCountries(ctx context.Context, _ string) ([]core.CatalogCountry, error) {
	countries, err := c.ListCountries(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]core.CatalogCountry, 0, len(countries))
	for _, country := range countries {
		iso2, callingCode := smsbowerCountryCodes(country.Name)
		out = append(out, core.CatalogCountry{
			CountryISO2:        iso2,
			Name:               country.Name,
			CountryCallingCode: callingCode,
		})
	}
	return out, nil
}
