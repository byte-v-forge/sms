package herosms

import (
	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/jsonx"
)

func activationOfferDefaultTier(offer activationOffer) []activationOfferPurchaseTier {
	price := jsonx.FirstScalar(offer.Prices.Min, offer.Prices.Default, offer.Prices.Retail)
	if price == "" {
		return nil
	}
	count := activationOfferDefaultCount(offer)
	if count <= 0 {
		return nil
	}
	return []activationOfferPurchaseTier{{
		Price:          core.Money{CurrencyCode: "USD", AmountDecimal: price},
		AvailableCount: count,
	}}
}

func activationOfferDefaultCount(offer activationOffer) int {
	if offer.Counts.DefaultPrice > 0 {
		return offer.Counts.DefaultPrice
	}
	return offer.Counts.Total
}
