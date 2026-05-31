package herosms

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
)

func (c *Client) ListNumberStatuses(ctx context.Context, countryID string) ([]PriceOffer, error) {
	params := url.Values{}
	if strings.TrimSpace(countryID) != "" {
		params.Set("country", strings.TrimSpace(countryID))
	}
	result, err := c.api.Do(ctx, "getNumbersStatus", params)
	if err != nil {
		return nil, err
	}
	var raw map[string]json.RawMessage
	if err := decodeHeroSMSJSONObject(result, &raw); err != nil {
		return nil, err
	}
	offers := make([]PriceOffer, 0, len(raw))
	for service, count := range raw {
		offers = append(offers, PriceOffer{
			CountryID:          strings.TrimSpace(countryID),
			UpstreamServiceKey: normalizeHeroSMSServiceKey(service),
			AvailableCount:     heroSMSInt(count),
		})
	}
	return offers, nil
}
