package app

import smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"

func filteredCatalogProviderConfigs(configs []*smsinternalv1.SmsProviderConfig, filter map[string]struct{}) []*smsinternalv1.SmsProviderConfig {
	if len(filter) == 0 {
		return configs
	}
	out := make([]*smsinternalv1.SmsProviderConfig, 0, len(configs))
	for _, config := range configs {
		if providerIncluded(config.GetProviderKey(), filter) {
			out = append(out, config)
		}
	}
	return out
}

func singleProviderKey(providerKeys []string) string {
	if len(providerKeys) == 1 {
		return providerKeys[0]
	}
	return ""
}
