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
	priceOffers, err := c.listOperatorRoutePriceOffers(ctx, query.ApplicationKey, country.ID)
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
		route := core.Route{
			ProviderKey:        ProviderKey,
			ApplicationKey:     offer.UpstreamServiceKey,
			UpstreamServiceKey: offer.UpstreamServiceKey,
			CountryISO2:        metadata.ISO2,
			CountryCallingCode: metadata.CallingCode,
			ProviderCountryID:  offer.CountryID,
			UpstreamProviderID: offer.Operator,
		}
		out = append(out, core.RouteOffer{
			ProviderKey:          ProviderKey,
			UpstreamProviderID:   offer.Operator,
			UpstreamProviderName: offer.Operator,
			ApplicationKey:       offer.UpstreamServiceKey,
			ApplicationName:      heroSMSApplicationName(offer.UpstreamServiceKey),
			CountryISO2:          metadata.ISO2,
			CountryName:          firstHeroSMSString(metadata.Name, offer.CountryID),
			CountryCallingCode:   metadata.CallingCode,
			Price:                offer.Price,
			AvailableCount:       offer.AvailableCount,
			ObservedAt:           now,
			Route:                route,
		})
	}
	return out, nil
}

func (c *Client) listRoutePriceOffers(ctx context.Context, applicationKey, countryID, operator string) ([]PriceOffer, error) {
	candidates := heroSMSServiceCandidates(applicationKey)
	var out []PriceOffer
	var lastErr error
	for _, service := range candidates {
		offers, err := c.ListPriceOffers(ctx, service, countryID, operator)
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
	statuses, err := c.ListNumberStatuses(ctx, countryID)
	if err != nil {
		return nil, lastErr
	}
	return uniqueHeroSMSOffers(statuses), nil
}
