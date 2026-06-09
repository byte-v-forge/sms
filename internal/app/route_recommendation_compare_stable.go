package app

func stableRouteCandidateLess(left routeCandidate, right routeCandidate) bool {
	return routeCandidateKey(left) < routeCandidateKey(right)
}
