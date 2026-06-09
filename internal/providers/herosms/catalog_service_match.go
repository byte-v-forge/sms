package herosms

func heroSMSApplicationMatches(serviceKey, queryApplicationKey string) bool {
	query := normalizeHeroSMSCatalogToken(queryApplicationKey)
	if query == "" {
		return true
	}
	service := normalizeHeroSMSCatalogToken(serviceKey)
	return service == query
}
