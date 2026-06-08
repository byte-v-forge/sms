package herosms

import "context"

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
