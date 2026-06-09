package app

import "github.com/byte-v-forge/sms/internal/core"

func withinPriceRange(price core.Money, minPrice core.Money, maxPrice core.Money) bool {
	offerAmount, ok := parseDecimalAmount(price.AmountDecimal)
	if !ok {
		return !moneyIsSet(minPrice) && !moneyIsSet(maxPrice)
	}
	if !withinPriceBound(price, offerAmount, minPrice, priceBoundMinimum) {
		return false
	}
	return withinPriceBound(price, offerAmount, maxPrice, priceBoundMaximum)
}
