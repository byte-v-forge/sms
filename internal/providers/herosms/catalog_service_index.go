package herosms

import (
	"strings"
)

func heroSMSServiceNameIndex(services []serviceMetadata) map[string]string {
	names := make(map[string]string, len(services))
	for _, service := range services {
		key := normalizeHeroSMSServiceKey(service.Service)
		name := strings.TrimSpace(service.Name)
		if key == "" || name == "" {
			continue
		}
		names[key] = name
	}
	return names
}
