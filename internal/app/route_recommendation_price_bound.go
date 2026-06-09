package app

import (
	"math/big"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

const (
	priceBoundMinimum = 1
	priceBoundMaximum = -1
)

func withinPriceBound(price core.Money, offerAmount *big.Rat, bound core.Money, direction int) bool {
	if strings.TrimSpace(bound.AmountDecimal) == "" {
		return true
	}
	if priceCurrencyMismatch(price, bound) {
		return false
	}
	boundAmount, ok := parseDecimalAmount(bound.AmountDecimal)
	if !ok {
		return false
	}
	return priceComparisonAllowed(offerAmount.Cmp(boundAmount), direction)
}

func priceComparisonAllowed(comparison int, direction int) bool {
	if direction > 0 {
		return comparison >= 0
	}
	return comparison <= 0
}
