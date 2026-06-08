package herosms

func heroSMSApplicationMatches(serviceKey, queryApplicationKey string) bool {
	query := normalizeHeroSMSCatalogToken(queryApplicationKey)
	if query == "" {
		return true
	}
	service := normalizeHeroSMSCatalogToken(serviceKey)
	if service == query {
		return true
	}
	for _, alias := range heroSMSServiceAliases[query] {
		if normalizeHeroSMSCatalogToken(alias) == service {
			return true
		}
	}
	return normalizeHeroSMSCatalogToken(heroSMSApplicationName(serviceKey)) == query
}
