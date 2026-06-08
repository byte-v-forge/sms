package app

import "sort"

func sortedProviderFilterKeys(providerFilter map[string]struct{}) []string {
	keys := make([]string, 0, len(providerFilter))
	for providerKey := range providerFilter {
		keys = append(keys, providerKey)
	}
	sort.Strings(keys)
	return keys
}

func normalizedProviderFilter(providerKeys []string) map[string]struct{} {
	if len(providerKeys) == 0 {
		return nil
	}
	filter := make(map[string]struct{}, len(providerKeys))
	for _, providerKey := range providerKeys {
		key := normalizeProviderKey(providerKey)
		if key == "" {
			continue
		}
		filter[key] = struct{}{}
	}
	if len(filter) == 0 {
		return nil
	}
	return filter
}
