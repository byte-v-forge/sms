package herosms

import (
	"strconv"
	"strings"
)

func heroSMSServiceCandidates(applicationKey string) []string {
	normalized := normalizeHeroSMSCatalogToken(applicationKey)
	if normalized == "" {
		return []string{""}
	}
	candidates := []string{strings.TrimSpace(applicationKey)}
	candidates = append(candidates, heroSMSServiceAliases[normalized]...)
	return uniqueHeroSMSStrings(candidates)
}

func heroSMSApplicationMatches(serviceKey, queryApplicationKey string) bool {
	query := normalizeHeroSMSCatalogToken(queryApplicationKey)
	if query == "" {
		return true
	}
	service := normalizeHeroSMSCatalogToken(serviceKey)
	if service == query {
		return true
	}
	for _, alias := range heroSMSServiceAliases[query] {
		if normalizeHeroSMSCatalogToken(alias) == service {
			return true
		}
	}
	return normalizeHeroSMSCatalogToken(heroSMSApplicationName(serviceKey)) == query
}

func heroSMSApplicationName(serviceKey string) string {
	if name := heroSMSServiceNames[normalizeHeroSMSServiceKey(serviceKey)]; name != "" {
		return name
	}
	return strings.TrimSpace(serviceKey)
}

func uniqueHeroSMSOffers(offers []PriceOffer) []PriceOffer {
	seen := map[string]struct{}{}
	out := make([]PriceOffer, 0, len(offers))
	for _, offer := range offers {
		key := offer.CountryID + "\x00" + offer.UpstreamServiceKey + "\x00" + offer.Operator + "\x00" + offer.Price.AmountDecimal
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, offer)
	}
	return out
}

func normalizeHeroSMSServiceKey(value string) string {
	value = strings.TrimSpace(value)
	if index := strings.LastIndex(value, "_"); index > 0 {
		suffix := value[index+1:]
		if _, err := strconv.Atoi(suffix); err == nil {
			return value[:index]
		}
	}
	return value
}

func normalizeHeroSMSCatalogToken(value string) string {
	var out strings.Builder
	for _, r := range strings.ToLower(strings.TrimSpace(value)) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			out.WriteRune(r)
		}
	}
	return out.String()
}

func uniqueHeroSMSStrings(values []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if _, exists := seen[value]; value == "" || exists {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}

func firstHeroSMSString(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
