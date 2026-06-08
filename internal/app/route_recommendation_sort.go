package app

import (
	"sort"
	"strings"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
)

func sortRouteCandidates(candidates []routeCandidate, strategy smsv1.SmsRouteStrategy) {
	sort.SliceStable(candidates, func(i, j int) bool {
		left := candidates[i]
		right := candidates[j]
		switch strategy {
		case smsv1.SmsRouteStrategy_SMS_ROUTE_STRATEGY_LOWEST_PRICE:
			if less, ok := priceLess(left, right); ok {
				return less
			}
			if left.offer.AvailableCount != right.offer.AvailableCount {
				return left.offer.AvailableCount > right.offer.AvailableCount
			}
		case smsv1.SmsRouteStrategy_SMS_ROUTE_STRATEGY_MOST_AVAILABLE:
			if left.offer.AvailableCount != right.offer.AvailableCount {
				return left.offer.AvailableCount > right.offer.AvailableCount
			}
			if less, ok := priceLess(left, right); ok {
				return less
			}
		default:
			if left.score != right.score {
				return left.score > right.score
			}
			if less, ok := priceLess(left, right); ok {
				return less
			}
			if left.offer.AvailableCount != right.offer.AvailableCount {
				return left.offer.AvailableCount > right.offer.AvailableCount
			}
		}
		return routeCandidateKey(left) < routeCandidateKey(right)
	})
}

func priceLess(left routeCandidate, right routeCandidate) (bool, bool) {
	if left.hasPrice != right.hasPrice {
		return left.hasPrice, true
	}
	if left.hasPrice {
		if cmp := left.price.Cmp(right.price); cmp != 0 {
			return cmp < 0, true
		}
	}
	return false, false
}

func routeCandidateKey(candidate routeCandidate) string {
	offer := candidate.offer
	return strings.Join([]string{
		normalizeProviderKey(offer.ProviderKey),
		routeText(offer.ApplicationKey),
		routeCountryISO2(offer.CountryISO2),
		routeCallingCode(offer.CountryCallingCode),
		routeText(offer.UpstreamProviderID),
		routeText(offer.UpstreamProviderName),
	}, "\x00")
}
