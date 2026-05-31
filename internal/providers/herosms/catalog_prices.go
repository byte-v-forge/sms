package herosms

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func (c *Client) ListPriceOffers(ctx context.Context, serviceKey, countryID, operator string) ([]PriceOffer, error) {
	params := url.Values{}
	if strings.TrimSpace(serviceKey) != "" {
		params.Set("service", strings.TrimSpace(serviceKey))
	}
	if strings.TrimSpace(countryID) != "" {
		params.Set("country", strings.TrimSpace(countryID))
	}
	if strings.TrimSpace(operator) != "" {
		params.Set("operator", strings.TrimSpace(operator))
	}
	result, err := c.api.Do(ctx, "getPrices", params)
	if err != nil {
		return nil, err
	}
	var raw map[string]map[string]struct {
		Cost        json.RawMessage `json:"cost"`
		Price       json.RawMessage `json:"price"`
		RetailPrice json.RawMessage `json:"retail_price"`
		Count       json.RawMessage `json:"count"`
	}
	if err := decodeHeroSMSJSONObject(result, &raw); err != nil {
		return nil, err
	}
	offers := make([]PriceOffer, 0)
	for cID, byService := range raw {
		for service, offer := range byService {
			offers = append(offers, PriceOffer{
				CountryID:          strings.TrimSpace(cID),
				UpstreamServiceKey: normalizeHeroSMSServiceKey(service),
				Operator:           strings.TrimSpace(operator),
				Price:              core.Money{AmountDecimal: firstHeroSMSScalar(offer.Cost, offer.Price, offer.RetailPrice)},
				AvailableCount:     heroSMSInt(offer.Count),
			})
		}
	}
	return offers, nil
}
