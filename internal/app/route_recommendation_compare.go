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
