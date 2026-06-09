package herosms

import (
	"context"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func (c *Client) ListRouteOffers(ctx context.Context, query core.RouteOfferQuery) ([]core.RouteOffer, error) {
	services, err := c.ListServices(ctx)
	if err != nil {
		return nil, err
	}
	service := heroSMSServiceForQuery(query.ApplicationKey, services)
	if query.ApplicationKey != "" && service == "" {
		return nil, nil
	}
	countries, err := c.ListCountries(ctx)
	if err != nil {
		return nil, err
	}
	country := heroSMSCountryForQuery(query, countries)
	if (query.CountryISO2 != "" || query.CountryCallingCode != "") && country.ID == "" {
		return nil, nil
	}
	priceOffers, err := c.listRoutePriceOffers(ctx, service, country.ID)
	if err != nil {
		return nil, err
	}
	routeQuery := query
	if service != "" {
		routeQuery.ApplicationKey = service
	}
	return heroSMSRouteOffersFromPrices(routeQuery, priceOffers, countries, country, heroSMSServiceNameIndex(services), time.Now().UTC()), nil
}
