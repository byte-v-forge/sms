package geox

import (
	"strings"

	"github.com/biter777/countries"
)

func CountryRegionCode(value string) string {
	country := countryByAlpha(value)
	if !country.IsValid() {
		return ""
	}
	return regionShortCode(country.Region())
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
