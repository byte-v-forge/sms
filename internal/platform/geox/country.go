package geox

import (
	"regexp"
	"strings"
	"sync"

	"github.com/biter777/countries"
)

var countryIndex = struct {
	sync.Once
	byAlpha map[string]countries.CountryCode
}{}

var countryEmojiIndex = struct {
	sync.Once
	values []countryEmoji
}{}

var countryNameIndex = struct {
	sync.Once
	values []countryName
}{}

type countryEmoji struct {
	emoji  string
	alpha2 string
}

type countryName struct {
	name   string
	alpha2 string
}

func NormalizeCountryAlpha2(value string) string {
	country := countryByAlpha(value)
	if !country.IsValid() {
		return ""
	}
	return country.Alpha2()
}

func CountryAlpha2ByName(value string) string {
	return alpha2ByCountryName(value, false)
}

func CountryAlpha2InText(value string) string {
	return alpha2ByCountryName(value, true)
}

func CountryCallingCode(value string) string {
	country := countryByAlpha(value)
	if !country.IsValid() {
		country = countries.ByName(strings.TrimSpace(value))
	}
	if !country.IsValid() {
		return ""
	}
	for _, code := range country.CallCodes() {
		if code.IsValid() {
			return strings.TrimPrefix(code.String(), "+")
		}
	}
	return ""
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

func CountryRegionCode(value string) string {
	country := countryByAlpha(value)
	if !country.IsValid() {
		return ""
	}
	return regionShortCode(country.Region())
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

func NormalizeRegionCode(value string) string {
	value = strings.ToUpper(strings.TrimSpace(value))
	switch strings.ReplaceAll(value, "_", " ") {
	case "AF", "AFRICA":
		return "AF"
	case "NA", "NORTH AMERICA":
		return "NA"
	case "OC", "OCEANIA":
		return "OC"
	case "AN", "ANTARCTICA":
		return "AN"
	case "AS", "ASIA":
		return "AS"
	case "EU", "EUROPE":
		return "EU"
	case "SA", "SOUTH AMERICA":
		return "SA"
	default:
		return ""
	}
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

func countryByAlpha(value string) countries.CountryCode {
	value = strings.ToUpper(strings.TrimSpace(value))
	if value == "" {
		return countries.Unknown
	}
	countryIndex.Do(func() {
		countryIndex.byAlpha = make(map[string]countries.CountryCode, countries.Total()*2)
		for _, country := range countries.All() {
			if alpha2 := country.Alpha2(); alpha2 != "" {
				countryIndex.byAlpha[strings.ToUpper(alpha2)] = country
			}
			if alpha3 := country.Alpha3(); alpha3 != "" {
				countryIndex.byAlpha[strings.ToUpper(alpha3)] = country
			}
		}
	})
	return countryIndex.byAlpha[value]
}

func regionShortCode(region countries.RegionCode) string {
	switch region {
	case countries.RegionAF:
		return "AF"
	case countries.RegionNA:
		return "NA"
	case countries.RegionOC:
		return "OC"
	case countries.RegionAN:
		return "AN"
	case countries.RegionAS:
		return "AS"
	case countries.RegionEU:
		return "EU"
	case countries.RegionSA:
		return "SA"
	default:
		return ""
	}
}
