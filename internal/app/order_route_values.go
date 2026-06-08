package app

import "strings"

func routeText(value string) string {
	return strings.TrimSpace(value)
}

func routeCountryISO2(value string) string {
	return strings.ToUpper(routeText(value))
}

func routeCallingCode(value string) string {
	return strings.TrimPrefix(routeText(value), "+")
}
