package fivesim

import "github.com/byte-v-forge/sms/internal/core"

func fiveSimApplicationNameIndex(applications []core.CatalogApplication) map[string]string {
	names := make(map[string]string, len(applications))
	for _, app := range applications {
		if app.ApplicationKey != "" && app.DisplayName != "" {
			names[app.ApplicationKey] = app.DisplayName
		}
	}
	return names
}
