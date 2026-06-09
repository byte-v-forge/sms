package herosms

import "strings"

func heroSMSApplicationName(serviceKey string, names map[string]string) string {
	if name := names[normalizeHeroSMSServiceKey(serviceKey)]; name != "" {
		return name
	}
	return strings.TrimSpace(serviceKey)
}
