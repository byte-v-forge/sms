package fivesim

import (
	"context"

	"github.com/byte-v-forge/sms/internal/core"
)

func (c *Client) ListRouteOffers(ctx context.Context, query core.RouteOfferQuery) ([]core.RouteOffer, error) {
	countries, err := c.ListCountries(ctx)
	if err != nil {
		return nil, err
	}
	countryID := fiveSimCountryIDForQuery(query, countries)
	if (query.CountryISO2 != "" || query.CountryCallingCode != "") && countryID == "" {
		return nil, nil
	}
	priceOffers, err := c.ListPriceOffers(ctx, query.ApplicationKey, countryID)
	if err != nil {
		return nil, err
	}
	return fiveSimRouteOffers(priceOffers, fiveSimCountriesByID(countries), c.catalogApplicationNames(ctx)), nil
}
