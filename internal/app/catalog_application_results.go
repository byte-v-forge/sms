package app

import (
	"sort"

	"github.com/byte-v-forge/sms/internal/core"
)

func collectCatalogProviderApplications(results []catalogProviderApplicationsResult) ([]core.CatalogApplication, error) {
	items := map[string]core.CatalogApplication{}
	lastErr := applicationProviderError(results)
	for _, result := range results {
		for _, app := range result.applications {
			key := routeText(app.ApplicationKey)
			if key == "" {
				continue
			}
			items[key] = bestCatalogApplication(items[key], core.CatalogApplication{
				ApplicationKey: key,
				DisplayName:    firstNonEmpty(routeText(app.DisplayName), key),
			})
		}
	}
	applications := catalogApplicationValues(items)
	if len(applications) == 0 && lastErr != nil {
		return nil, lastErr
	}
	return applications, nil
}

func applicationProviderError(results []catalogProviderApplicationsResult) error {
	var lastErr error
	for _, result := range results {
		if result.err == nil {
			continue
		}
		lastErr = result.err
	}
	return lastErr
}

func bestCatalogApplication(left core.CatalogApplication, right core.CatalogApplication) core.CatalogApplication {
	if left.ApplicationKey == "" || len(right.DisplayName) > len(left.DisplayName) {
		return right
	}
	return left
}

func catalogApplicationValues(items map[string]core.CatalogApplication) []core.CatalogApplication {
	out := make([]core.CatalogApplication, 0, len(items))
	for _, item := range items {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].DisplayName < out[j].DisplayName })
	return out
}
