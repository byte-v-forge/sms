package smsbower

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/searchx"
)

func (c *Client) ListCatalogApplications(ctx context.Context, query core.CatalogApplicationQuery) ([]core.CatalogApplication, error) {
	applications, err := c.ListApplications(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]core.CatalogApplication, 0, len(applications))
	for _, app := range applications {
		application := core.CatalogApplication{
			ApplicationKey: app.ApplicationKey,
			DisplayName:    app.DisplayName,
			Aliases:        []string{app.ApplicationKey, app.UpstreamServiceKey, app.DisplayName},
		}
		if !smsbowerCatalogApplicationMatches(application, query.SearchText) {
			continue
		}
		out = append(out, application)
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

func smsbowerCatalogApplicationMatches(application core.CatalogApplication, searchText string) bool {
	if searchx.Token(searchText) == "" {
		return true
	}
	for _, value := range append([]string{application.ApplicationKey, application.DisplayName}, application.Aliases...) {
		if searchx.ContainsToken(value, searchText) {
			return true
		}
	}
	return false
}
