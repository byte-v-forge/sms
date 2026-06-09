package herosms

import "context"

func (c *Client) listRouteCountries(ctx context.Context, serviceKey string) ([]countryMetadata, error) {
	serviceKey = normalizeHeroSMSServiceKey(serviceKey)
	if serviceKey == "" {
		return c.ListCountries(ctx)
	}
	return c.ListServiceCountries(ctx, serviceKey)
}
