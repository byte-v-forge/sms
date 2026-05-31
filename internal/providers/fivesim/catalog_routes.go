package fivesim

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
	countryID := fiveSimCountryIDForQuery(query, countries)
	if (query.CountryISO2 != "" || query.CountryCallingCode != "") && countryID == "" {
		return nil, nil
	}
	priceOffers, err := c.ListPriceOffers(ctx, query.ApplicationKey, countryID)
	if err != nil {
		return nil, err
	}
	countryByID := map[string]Country{}
	for _, country := range countries {
		countryByID[country.CountryID] = country
	}
	out := make([]core.RouteOffer, 0, len(priceOffers))
	now := time.Now().UTC()
	for _, offer := range priceOffers {
		country := countryByID[offer.CountryID]
		route := core.Route{
			ProviderKey:        ProviderKey,
			ApplicationKey:     offer.UpstreamServiceKey,
			UpstreamServiceKey: offer.UpstreamServiceKey,
			CountryISO2:        country.CountryISO2,
			CountryCallingCode: country.CountryCallingCode,
			ProviderCountryID:  offer.CountryID,
			UpstreamProviderID: offer.Operator,
		}
		out = append(out, core.RouteOffer{
			ProviderKey:          ProviderKey,
			UpstreamProviderID:   offer.Operator,
			UpstreamProviderName: offer.Operator,
			ApplicationKey:       offer.UpstreamServiceKey,
			ApplicationName:      offer.UpstreamServiceKey,
			CountryISO2:          country.CountryISO2,
			CountryName:          country.Name,
			CountryCallingCode:   country.CountryCallingCode,
			Price:                offer.Price,
			AvailableCount:       offer.AvailableCount,
			ObservedAt:           now,
			Route:                route,
		})
	}
	return out, nil
}
