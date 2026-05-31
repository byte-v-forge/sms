package fivesim

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/byte-v-forge/sms/internal/core"
)

func (c *Client) ListPriceOffers(ctx context.Context, product, countryID string) ([]PriceOffer, error) {
	params := url.Values{}
	if product != "" {
		params.Set("product", product)
	}
	if countryID != "" {
		params.Set("country", countryID)
	}
	var raw map[string]map[string]map[string]struct {
		Cost  json.RawMessage `json:"cost"`
		Count int             `json:"count"`
		Rate  float64         `json:"rate"`
	}
	if err := c.getJSON(ctx, "/v1/guest/prices", params, false, &raw); err != nil {
		return nil, err
	}
	var offers []PriceOffer
	for country, byProduct := range raw {
		for productKey, byOperator := range byProduct {
			for operator, offer := range byOperator {
				offers = append(offers, PriceOffer{
					CountryID:          country,
					UpstreamServiceKey: productKey,
					Operator:           operator,
					Price:              core.Money{CurrencyCode: c.currencyCode, AmountDecimal: rawJSONScalar(offer.Cost)},
					AvailableCount:     offer.Count,
					SuccessRate:        offer.Rate,
				})
			}
		}
	}
	return offers, nil
}
