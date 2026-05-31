package smsbower

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/byte-v-forge/common-lib/stringx"
	"github.com/byte-v-forge/sms/internal/core"
)

func (c *Client) ListPriceOffers(ctx context.Context, serviceKey, countryID string) ([]PriceOffer, error) {
	params := url.Values{}
	params.Set("service", serviceKey)
	params.Set("country", countryID)
	result, err := c.api.Do(ctx, "getPricesV3", params)
	if err != nil {
		return nil, err
	}
	var raw map[string]map[string]map[string]struct {
		Count      int             `json:"count"`
		Price      json.RawMessage `json:"price"`
		ProviderID json.RawMessage `json:"provider_id"`
	}
	if err := decodeJSONObject(result, &raw); err != nil {
		return nil, err
	}
	var offers []PriceOffer
	for cID, byService := range raw {
		for svc, byProvider := range byService {
			for providerID, offer := range byProvider {
				offers = append(offers, PriceOffer{
					CountryID:          cID,
					UpstreamServiceKey: svc,
					ProviderID:         stringx.FirstNonEmpty(rawJSONScalar(offer.ProviderID), providerID),
					Price:              core.Money{AmountDecimal: rawJSONScalar(offer.Price)},
					AvailableCount:     offer.Count,
				})
			}
		}
	}
	return offers, nil
}
