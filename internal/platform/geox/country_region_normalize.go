package geox

import "strings"

var normalizedRegionCodes = map[string]string{
	"AF":            "AF",
	"AFRICA":        "AF",
	"NA":            "NA",
	"NORTH AMERICA": "NA",
	"OC":            "OC",
	"OCEANIA":       "OC",
	"AN":            "AN",
	"ANTARCTICA":    "AN",
	"AS":            "AS",
	"ASIA":          "AS",
	"EU":            "EU",
	"EUROPE":        "EU",
	"SA":            "SA",
	"SOUTH AMERICA": "SA",
}

func NormalizeRegionCode(value string) string {
	value = strings.ToUpper(strings.TrimSpace(value))
	return normalizedRegionCodes[strings.ReplaceAll(value, "_", " ")]
}
