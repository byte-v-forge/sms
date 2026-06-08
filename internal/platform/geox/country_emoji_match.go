package geox

import "strings"

func countryCodesFromEmoji(value string) []string {
	out := []string{}
	for _, country := range countryEmojis() {
		if strings.Contains(value, country.emoji) {
			out = append(out, country.alpha2)
		}
	}
	return out
}
