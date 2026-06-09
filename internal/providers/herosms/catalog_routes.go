package herosms

import (
	"context"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func (c *Client) ListRouteOffers(ctx context.Context, query core.RouteOfferQuery) ([]core.RouteOffer, error) {
	service, err := c.routeService(ctx, query.ApplicationKey)
	if err != nil {
		return nil, err
	}
	if service.Key == "" {
		return nil, nil
	}
	countries, err := c.listRouteCountries(ctx, service.Key)
	if err != nil {
		return nil, err
	}
	country := heroSMSCountryForQuery(query, countries)
	if (query.CountryISO2 != "" || query.CountryCallingCode != "") && country.ID == "" {
		return nil, nil
	}
	priceOffers, err := c.listRoutePriceOffers(ctx, service.Key, country.ID)
	if err != nil {
		return nil, err
	}
	routeQuery := query
	if service.Key != "" {
		routeQuery.ApplicationKey = service.Key
	}
	return heroSMSRouteOffersFromPrices(routeQuery, priceOffers, countries, country, service.NameByID, time.Now().UTC()), nil
}
