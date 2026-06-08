package app

func routeCandidateMaxAvailable(candidates []routeCandidate) int {
	maxAvailable := 0
	for _, candidate := range candidates {
		if candidate.offer.AvailableCount > maxAvailable {
			maxAvailable = candidate.offer.AvailableCount
		}
	}
	return maxAvailable
}

func normalizedAvailableScore(available int, maxAvailable int) int32 {
	if available <= 0 || maxAvailable <= 0 {
		return 0
	}
	return int32((available * 1000) / maxAvailable)
}
