package fivesim

import (
	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/searchx"
)

func fiveSimProductForQuery(applicationKey string, applications []core.CatalogApplication) string {
	return searchx.MatchKey(applicationKey, fiveSimProductCandidates(applications))
}

func fiveSimProductCandidates(applications []core.CatalogApplication) []searchx.Candidate {
	candidates := make([]searchx.Candidate, 0, len(applications))
	for _, app := range applications {
		candidates = append(candidates, searchx.Candidate{
			Key:  app.ApplicationKey,
			Name: app.DisplayName,
		})
	}
	return candidates
}
