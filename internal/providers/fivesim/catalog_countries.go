package fivesim

import "context"

func (c *Client) ListCountries(ctx context.Context) ([]Country, error) {
	var raw map[string]countryShape
	if err := c.getJSON(ctx, "/v1/guest/countries", nil, false, &raw); err != nil {
		return nil, err
	}
	return fiveSimCountriesFromShapes(raw), nil
}
