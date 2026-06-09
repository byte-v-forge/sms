package herosms

import "github.com/byte-v-forge/sms/internal/platform/searchx"

func heroSMSServiceForQuery(applicationKey string, services []serviceMetadata) string {
	return searchx.MatchKey(applicationKey, heroSMSServiceSearchCandidates(services))
}

func heroSMSServiceSearchCandidates(services []serviceMetadata) []searchx.Candidate {
	candidates := make([]searchx.Candidate, 0, len(services))
	for _, service := range services {
		candidates = append(candidates, searchx.Candidate{
			Key:  normalizeHeroSMSServiceKey(service.Service),
			Name: service.Name,
		})
	}
	return candidates
}
