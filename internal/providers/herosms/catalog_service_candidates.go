package herosms

import "strings"

func heroSMSServiceCandidates(applicationKey string) []string {
	normalized := normalizeHeroSMSCatalogToken(applicationKey)
	if normalized == "" {
		return []string{""}
	}
	if aliases := heroSMSServiceAliases[normalized]; len(aliases) > 0 {
		return uniqueHeroSMSStrings(aliases)
	}
	candidates := []string{strings.TrimSpace(applicationKey)}
	return uniqueHeroSMSStrings(candidates)
}
