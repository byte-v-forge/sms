package app

import (
	"math/big"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func withinPriceRange(price core.Money, minPrice core.Money, maxPrice core.Money) bool {
	offerAmount, ok := parseDecimalAmount(price.AmountDecimal)
	if !ok {
		return !moneyIsSet(minPrice) && !moneyIsSet(maxPrice)
	}
	if !withinPriceBound(price, offerAmount, minPrice, 1) {
		return false
	}
	return withinPriceBound(price, offerAmount, maxPrice, -1)
}

func withinPriceBound(price core.Money, offerAmount *big.Rat, bound core.Money, direction int) bool {
	if strings.TrimSpace(bound.AmountDecimal) == "" {
		return true
	}
	if bound.CurrencyCode != "" && price.CurrencyCode != "" && !strings.EqualFold(bound.CurrencyCode, price.CurrencyCode) {
		return false
	}
	boundAmount, ok := parseDecimalAmount(bound.AmountDecimal)
	if !ok {
		return false
	}
	comparison := offerAmount.Cmp(boundAmount)
	if direction > 0 {
		return comparison >= 0
	}
	return comparison <= 0
}

func parseDecimalAmount(value string) (*big.Rat, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, false
	}
	amount, ok := new(big.Rat).SetString(value)
	if !ok {
		return nil, false
	}
	return amount, true
}
