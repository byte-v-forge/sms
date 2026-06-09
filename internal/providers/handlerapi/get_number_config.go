package handlerapi

import "strings"

func normalizeGetNumberV2Config(config GetNumberV2Config) GetNumberV2Config {
	config.ProviderName = defaultGetNumberV2Text(config.ProviderName, "sms provider")
	config.CountryLabel = defaultGetNumberV2Text(config.CountryLabel, "country")
	config.ProviderIDLabel = defaultGetNumberV2Text(config.ProviderIDLabel, "upstream provider id")
	return config
}

func defaultGetNumberV2Text(value string, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}
