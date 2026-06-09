package herosms

import "context"

func (c *Client) ListCountries(ctx context.Context) ([]countryMetadata, error) {
	result, err := c.api.Do(ctx, "getCountries", nil)
	if err != nil {
		return nil, err
	}
	return decodeHeroSMSCountries(result)
}
