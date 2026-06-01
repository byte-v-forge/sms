package herosms

import (
	"context"
	"encoding/json"
	"math/big"
	"net/url"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

type activationOffersResponse struct {
	Data map[string]map[string]activationOffer `json:"data"`
}

type activationOffer struct {
	Prices struct {
		Default json.RawMessage `json:"default"`
		Retail  json.RawMessage `json:"retail"`
		Min     json.RawMessage `json:"min"`
	} `json:"prices"`
	Counts struct {
		Total        int `json:"total"`
		Physical     int `json:"physical"`
		DefaultPrice int `json:"defaultPrice"`
	} `json:"counts"`
	PriceMap map[string]int `json:"map"`
}

func (c *Client) ListPriceOffers(ctx context.Context, serviceKey, countryID string) ([]PriceOffer, error) {
	params := url.Values{}
	services := heroSMSServiceCandidates(serviceKey)
	if len(services) > 0 && services[0] != "" {
		params.Set("services", strings.Join(services, ","))
	}
	if strings.TrimSpace(countryID) != "" {
		params.Set("countries", strings.TrimSpace(countryID))
	}

	var response activationOffersResponse
	if err := c.getOpenAPIJSON(ctx, "/activations/offers", params, &response); err != nil {
		return nil, err
	}
	offers := make([]PriceOffer, 0)
	for service, byCountry := range response.Data {
		for cID, item := range byCountry {
			price, availableCount, ok := activationOfferPurchasePrice(item)
			if !ok {
				continue
			}
			offers = append(offers, PriceOffer{
				CountryID:          strings.TrimSpace(cID),
				UpstreamServiceKey: normalizeHeroSMSServiceKey(service),
				Price:              price,
				AvailableCount:     availableCount,
			})
		}
	}
	return offers, nil
}

func activationOfferPurchasePrice(offer activationOffer) (core.Money, int, bool) {
	if price, count, ok := lowestActivationPrice(offer.PriceMap); ok {
		return core.Money{AmountDecimal: price}, count, true
	}
	price := firstHeroSMSScalar(offer.Prices.Min, offer.Prices.Default, offer.Prices.Retail)
	if price == "" {
		return core.Money{}, 0, false
	}
	count := offer.Counts.DefaultPrice
	if count <= 0 {
		count = offer.Counts.Total
	}
	return core.Money{AmountDecimal: price}, count, count > 0
}

func lowestActivationPrice(priceMap map[string]int) (string, int, bool) {
	var bestText string
	var bestValue *big.Rat
	var bestCount int
	for text, count := range priceMap {
		if count <= 0 {
			continue
		}
		value, ok := new(big.Rat).SetString(strings.TrimSpace(text))
		if !ok {
			continue
		}
		if bestValue == nil || value.Cmp(bestValue) < 0 {
			bestText = strings.TrimSpace(text)
			bestValue = value
			bestCount = count
		}
	}
	return bestText, bestCount, bestValue != nil
}
