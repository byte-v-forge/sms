package app

import (
	"sort"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
)

func sortRouteCandidates(candidates []routeCandidate, strategy smsv1.SmsRouteStrategy) {
	less := routeCandidateLess(strategy)
	sort.SliceStable(candidates, func(i, j int) bool { return less(candidates[i], candidates[j]) })
}
