package app

import (
	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/platform/pagex"
)

func recommendationStrategy(policy *smsv1.SmsRoutePolicy) smsv1.SmsRouteStrategy {
	switch policy.GetStrategy() {
	case smsv1.SmsRouteStrategy_SMS_ROUTE_STRATEGY_LOWEST_PRICE:
		return smsv1.SmsRouteStrategy_SMS_ROUTE_STRATEGY_LOWEST_PRICE
	case smsv1.SmsRouteStrategy_SMS_ROUTE_STRATEGY_MOST_AVAILABLE:
		return smsv1.SmsRouteStrategy_SMS_ROUTE_STRATEGY_MOST_AVAILABLE
	default:
		return smsv1.SmsRouteStrategy_SMS_ROUTE_STRATEGY_LOWEST_PRICE
	}
}

func recommendationLimit(policy *smsv1.SmsRoutePolicy) int {
	return pagex.NormalizeLimit(int(policy.GetLimit()), defaultRouteRecommendationLimit, maxRouteRecommendationLimit)
}
