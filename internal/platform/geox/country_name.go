package geox

import (
	"regexp"
	"strings"
	"sync"

	"github.com/biter777/countries"
)

var countryNameIndex = struct {
	sync.Once
	values []countryName
}{}

type countryName struct {
	name   string
	alpha2 string
}

func alpha2ByCountryName(value string, inText bool) string {
	normalized := normalizeCountryNameText(value)
	if normalized == "" {
		return ""
	}
	country := countries.ByName(normalized)
	if country.IsValid() {
		return country.Alpha2()
	}
	countryNameIndex.Do(func() {
		countryNameIndex.values = make([]countryName, 0, countries.Total())
		for _, country := range countries.All() {
			if !country.IsValid() {
				continue
			}
			countryNameIndex.values = append(countryNameIndex.values, countryName{name: normalizeCountryNameText(country.String()), alpha2: country.Alpha2()})
		}
	})
	for _, country := range countryNameIndex.values {
		if country.name == "" {
			continue
		}
		if normalized == country.name || inText && countryNameOccursInText(normalized, country.name) {
			return country.alpha2
		}
	}
	return ""
}

func countryNameOccursInText(text, name string) bool {
	text = " " + text + " "
	name = " " + name + " "
	if strings.Contains(text, name) {
		return true
	}
	return strings.Contains(text, strings.TrimSuffix(name, " ")+"s ")
}

var countryNameSeparator = regexp.MustCompile(`[^a-z0-9]+`)

func normalizeCountryNameText(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = countryNameSeparator.ReplaceAllString(value, " ")
	return strings.Join(strings.Fields(value), " ")
}
