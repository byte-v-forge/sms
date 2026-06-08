package smsbower

import (
	"context"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

type routeCatalogInputs struct {
	applications []ApplicationOffer
	countries    []Country
	serviceKey   string
	countryID    string
	empty        bool
}

func (c *Client) routeCatalogInputs(ctx context.Context, query core.RouteOfferQuery) (routeCatalogInputs, error) {
	applications, err := c.ListApplications(ctx)
	if err != nil {
		return routeCatalogInputs{}, err
	}
	countries, err := c.ListCountries(ctx)
	if err != nil {
		return routeCatalogInputs{}, err
	}
	service := matchService(query.ApplicationKey, applications)
	if strings.TrimSpace(query.ApplicationKey) != "" && service == "" {
		return routeCatalogInputs{empty: true}, nil
	}
	countryID := smsbowerCountryIDForQuery(query, countries)
	if (query.CountryISO2 != "" || query.CountryCallingCode != "") && countryID == "" {
		return routeCatalogInputs{empty: true}, nil
	}
	return routeCatalogInputs{applications: applications, countries: countries, serviceKey: service, countryID: countryID}, nil
}
