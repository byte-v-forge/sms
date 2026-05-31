package main

import smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"

func providerOptionJSON(plugin *smsinternalv1.SmsProviderPluginDescriptor) smsProviderOptionJSON {
	return smsProviderOptionJSON{
		ProviderKey: plugin.GetProviderKey(),
		DisplayName: plugin.GetDisplayName(),
	}
}

func providerSettingJSON(config *smsinternalv1.SmsProviderConfig) smsProviderSettingJSON {
	return smsProviderSettingJSON{
		ProviderKey: config.GetProviderKey(),
		Enabled:     config.GetEnabled(),
		APIKeySet:   config.GetCredentialSecretSet(),
	}
}
