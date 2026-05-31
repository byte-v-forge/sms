package smsbower

import (
	"context"
	"strings"
	"time"

	"github.com/byte-v-forge/common-lib/stringx"
	"github.com/byte-v-forge/sms/internal/core"
)

func (c *Client) ListRouteOffers(ctx context.Context, query core.RouteOfferQuery) ([]core.RouteOffer, error) {
	applications, err := c.ListApplications(ctx)
	if err != nil {
		return nil, err
	}
	countries, err := c.ListCountries(ctx)
	if err != nil {
		return nil, err
	}
	service := matchService(query.ApplicationKey, applications)
	if strings.TrimSpace(query.ApplicationKey) != "" && service == "" {
		return nil, nil
	}
	countryID := smsbowerCountryIDForQuery(query, countries)
	if (query.CountryISO2 != "" || query.CountryCallingCode != "") && countryID == "" {
		return nil, nil
	}
	priceOffers, err := c.ListPriceOffers(ctx, service, countryID)
	if err != nil {
		return nil, err
	}
	appNames := map[string]string{}
	for _, app := range applications {
		appNames[app.UpstreamServiceKey] = stringx.FirstNonEmpty(app.DisplayName, app.UpstreamServiceKey)
	}
	countryByID := map[string]Country{}
	for _, country := range countries {
		countryByID[country.CountryID] = country
	}
	out := make([]core.RouteOffer, 0, len(priceOffers))
	now := time.Now().UTC()
	for _, offer := range priceOffers {
		country := countryByID[offer.CountryID]
		iso2, callingCode := smsbowerCountryCodes(country.Name)
		route := core.Route{
			ProviderKey:        ProviderKey,
			ApplicationKey:     offer.UpstreamServiceKey,
			UpstreamServiceKey: offer.UpstreamServiceKey,
			CountryISO2:        iso2,
			CountryCallingCode: callingCode,
			ProviderCountryID:  offer.CountryID,
			UpstreamProviderID: offer.ProviderID,
		}
		out = append(out, core.RouteOffer{
			ProviderKey:          ProviderKey,
			UpstreamProviderID:   offer.ProviderID,
			UpstreamProviderName: offer.ProviderID,
			ApplicationKey:       offer.UpstreamServiceKey,
			ApplicationName:      stringx.FirstNonEmpty(appNames[offer.UpstreamServiceKey], offer.UpstreamServiceKey),
			CountryISO2:          iso2,
			CountryName:          stringx.FirstNonEmpty(country.Name, offer.CountryID),
			CountryCallingCode:   callingCode,
			Price:                offer.Price,
			AvailableCount:       offer.AvailableCount,
			ObservedAt:           now,
			Route:                route,
		})
	}
	return out, nil
}
