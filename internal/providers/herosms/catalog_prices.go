package herosms

import (
	"context"
	"net/url"
	"strings"
)

func (c *Client) ListPriceOffers(ctx context.Context, serviceKey, countryID string) ([]PriceOffer, error) {
	params := url.Values{}
	services := heroSMSServiceCandidates(serviceKey)
	if len(services) > 0 && services[0] != "" {
		params.Set("services", strings.Join(services, ","))
	}
	if strings.TrimSpace(countryID) != "" {
		params.Set("countries", strings.TrimSpace(countryID))
	}

	var response activationOffersResponse
	if err := c.getOpenAPIJSON(ctx, "/activations/offers", params, &response); err != nil {
		return nil, err
	}
	return activationOffers(response), nil
}

func activationOffers(response activationOffersResponse) []PriceOffer {
	offers := make([]PriceOffer, 0)
	for service, byCountry := range response.Data {
		for cID, item := range byCountry {
			for _, tier := range activationOfferPurchaseTiers(item) {
				offers = append(offers, PriceOffer{
					CountryID:          strings.TrimSpace(cID),
					UpstreamServiceKey: normalizeHeroSMSServiceKey(service),
					Price:              tier.Price,
					AvailableCount:     tier.AvailableCount,
				})
			}
		}
	}
	return offers
}
