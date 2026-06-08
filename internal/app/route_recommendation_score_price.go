package app

import "math/big"

func routeCandidatePriceRange(candidates []routeCandidate) (*big.Rat, *big.Rat) {
	var minPrice *big.Rat
	var maxPrice *big.Rat
	for _, candidate := range candidates {
		if !candidate.hasPrice {
			continue
		}
		if minPrice == nil || candidate.price.Cmp(minPrice) < 0 {
			minPrice = new(big.Rat).Set(candidate.price)
		}
		if maxPrice == nil || candidate.price.Cmp(maxPrice) > 0 {
			maxPrice = new(big.Rat).Set(candidate.price)
		}
	}
	return minPrice, maxPrice
}

func normalizedPriceScore(price *big.Rat, hasPrice bool, minPrice *big.Rat, maxPrice *big.Rat) int32 {
	if !hasPrice || minPrice == nil || maxPrice == nil {
		return 0
	}
	if maxPrice.Cmp(minPrice) == 0 {
		return 1000
	}
	numerator := new(big.Rat).Sub(maxPrice, price)
	denominator := new(big.Rat).Sub(maxPrice, minPrice)
	score := new(big.Rat).Quo(numerator, denominator)
	score.Mul(score, big.NewRat(1000, 1))
	return ratScore(score)
}
