package herosms

import (
	"math/big"
	"sort"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func activationOfferPriceMapTiers(priceMap map[string]int) []activationOfferPurchaseTier {
	tiers := make([]activationOfferPurchaseTier, 0, len(priceMap))
	for text, count := range priceMap {
		tier := activationOfferPriceMapTier(text, count)
		if tier == nil {
			continue
		}
		tiers = append(tiers, *tier)
	}
	sortActivationOfferPurchaseTiers(tiers)
	return tiers
}

func activationOfferPriceMapTier(text string, count int) *activationOfferPurchaseTier {
	if count <= 0 {
		return nil
	}
	priceText := strings.TrimSpace(text)
	value, ok := new(big.Rat).SetString(priceText)
	if !ok {
		return nil
	}
	return &activationOfferPurchaseTier{
		Price:          core.Money{CurrencyCode: "USD", AmountDecimal: priceText},
		AvailableCount: count,
		amount:         value,
	}
}

func sortActivationOfferPurchaseTiers(tiers []activationOfferPurchaseTier) {
	sort.SliceStable(tiers, func(i, j int) bool {
		return tiers[i].amount.Cmp(tiers[j].amount) < 0
	})
	for i := range tiers {
		tiers[i].amount = nil
	}
}
