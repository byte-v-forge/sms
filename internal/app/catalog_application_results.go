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
			normalized := normalizedCatalogApplication(app)
			identity := catalogApplicationIdentity(normalized)
			if identity == "" {
				continue
			}
			items[identity] = bestCatalogApplication(identity, items[identity], normalized)
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

func bestCatalogApplication(identity string, left core.CatalogApplication, right core.CatalogApplication) core.CatalogApplication {
	if left.ApplicationKey == "" {
		return right
	}
	display := left.DisplayName
	if len(right.DisplayName) > len(display) {
		display = right.DisplayName
	}
	return core.CatalogApplication{
		ApplicationKey: bestCatalogApplicationKey(identity, left.ApplicationKey, right.ApplicationKey, display),
		DisplayName:    display,
		Aliases:        uniqueCatalogValues(append(left.Aliases, right.Aliases...)),
	}
}

func bestCatalogApplicationKey(identity string, left string, right string, display string) string {
	if routeText(right) == identity {
		return right
	}
	if routeText(left) == identity {
		return left
	}
	if normalizeCatalogToken(left) == identity {
		return left
	}
	if normalizeCatalogToken(right) == identity {
		return right
	}
	return firstNonEmpty(display, left, right)
}

func catalogApplicationValues(items map[string]core.CatalogApplication) []core.CatalogApplication {
	out := make([]core.CatalogApplication, 0, len(items))
	for _, item := range items {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].DisplayName < out[j].DisplayName })
	return out
}
