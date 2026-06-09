package smsbower

import "github.com/byte-v-forge/sms/internal/platform/searchx"

func matchService(candidate string, applications []ApplicationOffer) string {
	return searchx.MatchKey(candidate, smsbowerServiceCandidates(applications))
}

func smsbowerServiceCandidates(applications []ApplicationOffer) []searchx.Candidate {
	candidates := make([]searchx.Candidate, 0, len(applications))
	for _, app := range applications {
		candidates = append(candidates, searchx.Candidate{
			Key:  app.UpstreamServiceKey,
			Name: app.DisplayName,
		})
	}
	return candidates
}
