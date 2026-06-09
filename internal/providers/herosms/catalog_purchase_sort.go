package herosms

import (
	"math/big"
	"strings"
)

func comparePriceOffers(left, right PriceOffer) int {
	if result := comparePriceText(left.Price.AmountDecimal, right.Price.AmountDecimal); result != 0 {
		return result
	}
	if result := strings.Compare(left.Operator, right.Operator); result != 0 {
		return result
	}
	return strings.Compare(left.UpstreamServiceKey, right.UpstreamServiceKey)
}

func comparePriceText(left, right string) int {
	leftValue, leftOK := new(big.Rat).SetString(strings.TrimSpace(left))
	rightValue, rightOK := new(big.Rat).SetString(strings.TrimSpace(right))
	if leftOK && rightOK {
		return leftValue.Cmp(rightValue)
	}
	return strings.Compare(strings.TrimSpace(left), strings.TrimSpace(right))
}
