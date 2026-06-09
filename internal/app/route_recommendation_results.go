package app

func routeRecommendations(candidates []routeCandidate, limit int) []RouteRecommendation {
	if limit > len(candidates) {
		limit = len(candidates)
	}
	recommendations := make([]RouteRecommendation, 0, limit)
	for index := 0; index < limit; index++ {
		recommendations = append(recommendations, RouteRecommendation{
			Offer: candidates[index].offer,
			Score: candidates[index].score,
		})
	}
	return recommendations
}
