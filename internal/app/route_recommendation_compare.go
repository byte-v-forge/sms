package app

import smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"

type routeCandidateComparator func(left routeCandidate, right routeCandidate) bool

func routeCandidateLess(strategy smsv1.SmsRouteStrategy) routeCandidateComparator {
	switch strategy {
	case smsv1.SmsRouteStrategy_SMS_ROUTE_STRATEGY_LOWEST_PRICE:
		return lowestPriceCandidateLess
	case smsv1.SmsRouteStrategy_SMS_ROUTE_STRATEGY_MOST_AVAILABLE:
		return mostAvailableCandidateLess
	default:
		return bestScoreCandidateLess
	}
}

func lowestPriceCandidateLess(left routeCandidate, right routeCandidate) bool {
	if less, ok := priceLess(left, right); ok {
		return less
	}
	if availableCountDiffers(left, right) {
		return moreAvailable(left, right)
	}
	return stableRouteCandidateLess(left, right)
}

func mostAvailableCandidateLess(left routeCandidate, right routeCandidate) bool {
	if availableCountDiffers(left, right) {
		return moreAvailable(left, right)
	}
	if less, ok := priceLess(left, right); ok {
		return less
	}
	return stableRouteCandidateLess(left, right)
}

func bestScoreCandidateLess(left routeCandidate, right routeCandidate) bool {
	if left.score != right.score {
		return left.score > right.score
	}
	if less, ok := priceLess(left, right); ok {
		return less
	}
	if availableCountDiffers(left, right) {
		return moreAvailable(left, right)
	}
	return stableRouteCandidateLess(left, right)
}

func availableCountDiffers(left routeCandidate, right routeCandidate) bool {
	return left.offer.AvailableCount != right.offer.AvailableCount
}

func moreAvailable(left routeCandidate, right routeCandidate) bool {
	return left.offer.AvailableCount > right.offer.AvailableCount
}

func stableRouteCandidateLess(left routeCandidate, right routeCandidate) bool {
	return routeCandidateKey(left) < routeCandidateKey(right)
}
