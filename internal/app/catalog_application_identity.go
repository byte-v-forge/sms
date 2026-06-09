package app

import (
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func normalizedCatalogApplication(app core.CatalogApplication) core.CatalogApplication {
	key := routeText(app.ApplicationKey)
	display := firstNonEmpty(routeText(app.DisplayName), key)
	return core.CatalogApplication{
		ApplicationKey: catalogApplicationQueryKey(key, display),
		DisplayName:    display,
		Aliases:        normalizedCatalogApplicationAliases(key, display, app.Aliases),
	}
}

func catalogApplicationIdentity(app core.CatalogApplication) string {
	if token := normalizeCatalogToken(catalogApplicationPrimaryName(app.DisplayName)); token != "" {
		return token
	}
	if token := normalizeCatalogToken(app.DisplayName); token != "" {
		return token
	}
	return normalizeCatalogToken(app.ApplicationKey)
}

func catalogApplicationPrimaryName(value string) string {
	for _, separator := range []string{",", "/", "|", ";"} {
		if index := strings.Index(value, separator); index >= 0 {
			return routeText(value[:index])
		}
	}
	return routeText(value)
}

func catalogApplicationQueryKey(key string, display string) string {
	return firstNonEmpty(key, display)
}

func normalizedCatalogApplicationAliases(key string, display string, aliases []string) []string {
	values := append([]string{key, display}, aliases...)
	return uniqueCatalogValues(values)
}

func uniqueCatalogValues(values []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = routeText(value)
		if value == "" {
			continue
		}
		token := normalizeCatalogToken(value)
		if token == "" {
			continue
		}
		if _, ok := seen[token]; ok {
			continue
		}
		seen[token] = struct{}{}
		out = append(out, value)
	}
	return out
}
