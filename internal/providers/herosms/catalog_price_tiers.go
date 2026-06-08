package herosms

import (
	"math/big"
	"sort"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/jsonx"
)

func activationOfferPurchaseTiers(offer activationOffer) []activationOfferPurchaseTier {
	if tiers := activationOfferPriceMapTiers(offer.PriceMap); len(tiers) > 0 {
		return tiers
	}
	price := jsonx.FirstScalar(offer.Prices.Min, offer.Prices.Default, offer.Prices.Retail)
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
