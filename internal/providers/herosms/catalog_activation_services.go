package herosms

import (
	"context"
	"sort"
)

func (c *Client) ListActivationServiceKeys(ctx context.Context) ([]string, error) {
	var response activationOffersResponse
	if err := c.getOpenAPIJSON(ctx, "/activations/offers", nil, &response); err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(response.Data))
	for service := range response.Data {
		if key := normalizeHeroSMSServiceKey(service); key != "" {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)
	return uniqueHeroSMSStrings(keys), nil
}
