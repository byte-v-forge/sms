package app

func mostAvailableCandidateLess(left routeCandidate, right routeCandidate) bool {
	if availableCountDiffers(left, right) {
		return moreAvailable(left, right)
	}
	if less, ok := priceLess(left, right); ok {
		return less
	}
	return stableRouteCandidateLess(left, right)
}

func availableCountDiffers(left routeCandidate, right routeCandidate) bool {
	return left.offer.AvailableCount != right.offer.AvailableCount
}

func moreAvailable(left routeCandidate, right routeCandidate) bool {
	return left.offer.AvailableCount > right.offer.AvailableCount
}
