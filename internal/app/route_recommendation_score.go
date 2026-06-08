package app

import smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"

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
