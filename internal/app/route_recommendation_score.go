package app

import (
	"math"
	"math/big"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
)

func scoreRouteCandidates(candidates []routeCandidate, strategy smsv1.SmsRouteStrategy) {
	minPrice, maxPrice := routeCandidatePriceRange(candidates)
	maxAvailable := routeCandidateMaxAvailable(candidates)
	for index := range candidates {
		priceScore := normalizedPriceScore(candidates[index].price, candidates[index].hasPrice, minPrice, maxPrice)
		availableScore := normalizedAvailableScore(candidates[index].offer.AvailableCount, maxAvailable)
		switch strategy {
		case smsv1.SmsRouteStrategy_SMS_ROUTE_STRATEGY_LOWEST_PRICE:
			candidates[index].score = priceScore
		case smsv1.SmsRouteStrategy_SMS_ROUTE_STRATEGY_MOST_AVAILABLE:
			candidates[index].score = availableScore
		default:
			candidates[index].score = int32((int(priceScore)*80 + int(availableScore)*20) / 100)
		}
	}
}

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

func routeCandidateMaxAvailable(candidates []routeCandidate) int {
	maxAvailable := 0
	for _, candidate := range candidates {
		if candidate.offer.AvailableCount > maxAvailable {
			maxAvailable = candidate.offer.AvailableCount
		}
	}
	return maxAvailable
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

func normalizedAvailableScore(available int, maxAvailable int) int32 {
	if available <= 0 || maxAvailable <= 0 {
		return 0
	}
	return int32((available * 1000) / maxAvailable)
}

func ratScore(value *big.Rat) int32 {
	if value == nil {
		return 0
	}
	score, _ := value.Float64()
	score = math.Round(score)
	if score < 0 {
		return 0
	}
	if score > 1000 {
		return 1000
	}
	return int32(score)
}
