package app

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
