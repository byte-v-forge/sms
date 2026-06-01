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
