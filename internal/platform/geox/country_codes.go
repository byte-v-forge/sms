package geox

import (
	"strings"
	"sync"

	"github.com/biter777/countries"
)

var countryEmojiIndex = struct {
	sync.Once
	values []countryEmoji
}{}

type countryEmoji struct {
	emoji  string
	alpha2 string
}

func CountryCodesInText(value string) []string {
	seen := map[string]struct{}{}
	out := []string{}
	addCountry := func(value string) {
		country := NormalizeCountryAlpha2(value)
		if country == "" {
			return
		}
		if _, exists := seen[country]; exists {
			return
		}
		seen[country] = struct{}{}
		out = append(out, country)
	}
	for _, value := range countryCodesFromEmoji(value) {
		addCountry(value)
	}
	for _, value := range alphaTokens(value) {
		addCountry(value)
	}
	return out
}

func countryCodesFromEmoji(value string) []string {
	countryEmojiIndex.Do(func() {
		countryEmojiIndex.values = make([]countryEmoji, 0, countries.Total())
		for _, country := range countries.All() {
			if !country.IsValid() {
				continue
			}
			countryEmojiIndex.values = append(countryEmojiIndex.values, countryEmoji{
				emoji:  country.Emoji(),
				alpha2: country.Alpha2(),
			})
		}
	})
	out := []string{}
	for _, country := range countryEmojiIndex.values {
		if strings.Contains(value, country.emoji) {
			out = append(out, country.alpha2)
		}
	}
	return out
}

func alphaTokens(value string) []string {
	fields := strings.FieldsFunc(value, func(r rune) bool {
		return !((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z'))
	})
	out := []string{}
	for _, field := range fields {
		field = strings.ToUpper(strings.TrimSpace(field))
		if len(field) == 2 || len(field) == 3 {
			out = append(out, field)
		}
	}
	return out
}
