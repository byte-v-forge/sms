package herosms

import (
	"context"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func (c *Client) ListRouteOffers(ctx context.Context, query core.RouteOfferQuery) ([]core.RouteOffer, error) {
	countries, err := c.ListCountries(ctx)
	if err != nil {
		return nil, err
	}
	country := heroSMSCountryForQuery(query, countries)
	if (query.CountryISO2 != "" || query.CountryCallingCode != "") && country.ID == "" {
		return nil, nil
	}
	priceOffers, err := c.listRoutePriceOffers(ctx, query.ApplicationKey, country.ID)
	if err != nil {
		return nil, err
	}
	return heroSMSRouteOffersFromPrices(query, priceOffers, countries, country, time.Now().UTC()), nil
}
