package herosms

import (
	"context"
	"net/url"
	"strings"
)

type operatorResponse struct {
	CountryOperators map[string][]string `json:"countryOperators"`
}

func (c *Client) ListOperators(ctx context.Context, countryID string) ([]string, error) {
	params := url.Values{}
	if strings.TrimSpace(countryID) != "" {
		params.Set("country", strings.TrimSpace(countryID))
	}
	result, err := c.api.Do(ctx, "getOperators", params)
	if err != nil {
		return nil, err
	}
	var response operatorResponse
	if err := decodeHeroSMSJSONObject(result, &response); err != nil {
		return nil, err
	}
	return uniqueHeroSMSStrings(response.CountryOperators[strings.TrimSpace(countryID)]), nil
}

func (c *Client) listOperatorRoutePriceOffers(ctx context.Context, applicationKey, countryID string) ([]PriceOffer, error) {
	operators, err := c.ListOperators(ctx, countryID)
	if err != nil {
		if !isHeroSMSUnsupportedCatalogLookup(err) {
			return nil, err
		}
		return c.listRoutePriceOffers(ctx, applicationKey, countryID, "")
	}
	if len(operators) == 0 {
		return c.listRoutePriceOffers(ctx, applicationKey, countryID, "")
	}

	var out []PriceOffer
	var lastErr error
	for _, operator := range operators {
		offers, err := c.listRoutePriceOffers(ctx, applicationKey, countryID, operator)
		if err != nil {
			lastErr = err
			continue
		}
		out = append(out, offers...)
	}
	if len(out) > 0 || lastErr == nil {
		return uniqueHeroSMSOffers(out), nil
	}
	return nil, lastErr
}
