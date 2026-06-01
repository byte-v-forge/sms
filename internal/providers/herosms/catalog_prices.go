package herosms

import (
	"context"
	"encoding/json"
	"math/big"
	"net/url"
	"sort"
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
			for _, tier := range activationOfferPurchaseTiers(item) {
				offers = append(offers, PriceOffer{
					CountryID:          strings.TrimSpace(cID),
					UpstreamServiceKey: normalizeHeroSMSServiceKey(service),
					Price:              tier.Price,
					AvailableCount:     tier.AvailableCount,
				})
			}
		}
	}
	return offers, nil
}

type activationOfferPurchaseTier struct {
	Price          core.Money
	AvailableCount int
	amount         *big.Rat
}

func activationOfferPurchaseTiers(offer activationOffer) []activationOfferPurchaseTier {
	if tiers := activationOfferPriceMapTiers(offer.PriceMap); len(tiers) > 0 {
		return tiers
	}
	price := firstHeroSMSScalar(offer.Prices.Min, offer.Prices.Default, offer.Prices.Retail)
	if price == "" {
		return nil
	}
	count := offer.Counts.DefaultPrice
	if count <= 0 {
		count = offer.Counts.Total
	}
	if count <= 0 {
		return nil
	}
	return []activationOfferPurchaseTier{
		{
			Price:          core.Money{CurrencyCode: "USD", AmountDecimal: price},
			AvailableCount: count,
		},
	}
}

func activationOfferPriceMapTiers(priceMap map[string]int) []activationOfferPurchaseTier {
	tiers := make([]activationOfferPurchaseTier, 0, len(priceMap))
	for text, count := range priceMap {
		if count <= 0 {
			continue
		}
		priceText := strings.TrimSpace(text)
		value, ok := new(big.Rat).SetString(priceText)
		if !ok {
			continue
		}
		tiers = append(tiers, activationOfferPurchaseTier{
			Price:          core.Money{CurrencyCode: "USD", AmountDecimal: priceText},
			AvailableCount: count,
			amount:         value,
		})
	}
	sort.SliceStable(tiers, func(i, j int) bool {
		return tiers[i].amount.Cmp(tiers[j].amount) < 0
	})
	for i := range tiers {
		tiers[i].amount = nil
	}
	return tiers
}
