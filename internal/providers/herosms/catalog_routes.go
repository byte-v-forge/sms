package herosms

import (
	"context"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/stringx"
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
	out := make([]core.RouteOffer, 0, len(priceOffers))
	now := time.Now().UTC()
	for _, offer := range priceOffers {
		if !heroSMSApplicationMatches(offer.UpstreamServiceKey, query.ApplicationKey) {
			continue
		}
		metadata := country
		if metadata.ID == "" {
			metadata = heroSMSCountryByID(countries, offer.CountryID)
		}
		applicationKey := stringx.FirstNonEmpty(query.ApplicationKey, heroSMSPublicApplicationKey(offer.UpstreamServiceKey), offer.UpstreamServiceKey)
		route := core.Route{
			ProviderKey:        ProviderKey,
			ApplicationKey:     applicationKey,
			UpstreamServiceKey: offer.UpstreamServiceKey,
			CountryISO2:        metadata.ISO2,
			CountryCallingCode: metadata.CallingCode,
			ProviderCountryID:  offer.CountryID,
			MaxPrice:           offer.Price,
		}
		out = append(out, core.RouteOffer{
			ProviderKey:          ProviderKey,
			UpstreamProviderName: "any",
			ApplicationKey:       applicationKey,
			ApplicationName:      heroSMSApplicationName(offer.UpstreamServiceKey),
			CountryISO2:          metadata.ISO2,
			CountryName:          stringx.FirstNonEmpty(metadata.Name, offer.CountryID),
			CountryCallingCode:   metadata.CallingCode,
			Price:                offer.Price,
			AvailableCount:       offer.AvailableCount,
			ObservedAt:           now,
			Route:                route,
		})
	}
	return out, nil
}

func (c *Client) listRoutePriceOffers(ctx context.Context, applicationKey, countryID string) ([]PriceOffer, error) {
	candidates := heroSMSServiceCandidates(applicationKey)
	var out []PriceOffer
	var lastErr error
	for _, service := range candidates {
		offers, err := c.ListPriceOffers(ctx, service, countryID)
		if err != nil {
			if isHeroSMSUnsupportedCatalogLookup(err) {
				lastErr = err
				continue
			}
			return nil, err
		}
		out = append(out, offers...)
	}
	if len(out) > 0 || lastErr == nil {
		return uniqueHeroSMSOffers(out), nil
	}
	return nil, lastErr
}
