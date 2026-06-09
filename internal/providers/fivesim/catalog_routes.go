package fivesim

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (c *Client) ListRouteOffers(ctx context.Context, query core.RouteOfferQuery) ([]core.RouteOffer, error) {
	applications, err := c.ListCatalogApplications(ctx)
	if err != nil {
		return nil, err
	}
	product := fiveSimProductForQuery(query.ApplicationKey, applications)
	if query.ApplicationKey != "" && product == "" {
		return nil, nil
	}
	countries, err := c.ListCountries(ctx)
	if err != nil {
		return nil, err
	}
	countryID := fiveSimCountryIDForQuery(query, countries)
	if (query.CountryISO2 != "" || query.CountryCallingCode != "") && countryID == "" {
		return nil, nil
	}
	priceOffers, err := c.ListPriceOffers(ctx, product, countryID)
	if err != nil {
		return nil, err
	}
	return fiveSimRouteOffers(priceOffers, fiveSimCountriesByID(countries), fiveSimApplicationNameIndex(applications)), nil
}
