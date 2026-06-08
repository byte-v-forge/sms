package herosms

import "strings"

func heroSMSApplicationName(serviceKey string) string {
	if name := heroSMSServiceNames[normalizeHeroSMSServiceKey(serviceKey)]; name != "" {
		return name
	}
	return strings.TrimSpace(serviceKey)
}

func heroSMSPublicApplicationKey(serviceKey string) string {
	switch normalizeHeroSMSServiceKey(serviceKey) {
	case "ni":
		return "gojek"
	case "wa":
		return "whatsapp"
	default:
		return strings.TrimSpace(serviceKey)
	}
}
