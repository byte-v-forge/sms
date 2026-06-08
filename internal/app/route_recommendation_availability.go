package app

import (
	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
)

func routeMinimumAvailability(policy *smsv1.SmsRoutePolicy) int {
	if policy.GetMinAvailableCount() > 0 {
		return int(policy.GetMinAvailableCount())
	}
	return 1
}

func routeCandidatesWithMinimumAvailability(candidates []routeCandidate, minimum int) []routeCandidate {
	if minimum <= 1 {
		return candidates
	}
	out := make([]routeCandidate, 0, len(candidates))
	for _, candidate := range candidates {
		if candidate.offer.AvailableCount >= minimum {
			out = append(out, candidate)
		}
	}
	return out
}
