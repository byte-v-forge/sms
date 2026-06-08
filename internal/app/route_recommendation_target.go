package app

import "github.com/byte-v-forge/sms/internal/core"

func normalizeRecommendationTarget(target core.Target) core.Target {
	target.ApplicationKey = routeText(target.ApplicationKey)
	target.CountryISO2 = routeCountryISO2(target.CountryISO2)
	target.CountryCallingCode = routeCallingCode(target.CountryCallingCode)
	return target
}

func validateRecommendationTarget(target core.Target) error {
	if target.ApplicationKey == "" {
		return core.NewError(core.CodeValidationFailed, "sms route recommendation target application_key is required", false)
	}
	if target.CountryISO2 == "" && target.CountryCallingCode == "" {
		return core.NewError(core.CodeValidationFailed, "sms route recommendation target country is required", false)
	}
	return nil
}
