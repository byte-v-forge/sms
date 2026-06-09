package app

func lowestPriceCandidateLess(left routeCandidate, right routeCandidate) bool {
	if less, ok := priceLess(left, right); ok {
		return less
	}
	if availableCountDiffers(left, right) {
		return moreAvailable(left, right)
	}
	return stableRouteCandidateLess(left, right)
}
