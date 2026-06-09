package app

func providerIncluded(providerKey string, filter map[string]struct{}) bool {
	if len(filter) == 0 {
		return true
	}
	_, ok := filter[normalizeProviderKey(providerKey)]
	return ok
}
