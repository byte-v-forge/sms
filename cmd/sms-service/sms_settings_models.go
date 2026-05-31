package main

type smsProviderOptionJSON struct {
	ProviderKey string `json:"provider_key"`
	DisplayName string `json:"display_name"`
}

type smsProviderSettingJSON struct {
	ProviderKey string `json:"provider_key"`
	Enabled     bool   `json:"enabled"`
	APIKeySet   bool   `json:"api_key_set"`
}

type listSMSProviderSettingsResponseJSON struct {
	ProviderOptions []smsProviderOptionJSON  `json:"provider_options"`
	Providers       []smsProviderSettingJSON `json:"providers"`
}

type saveSMSProviderSettingRequestJSON struct {
	ProviderKey string `json:"provider_key"`
	Enabled     *bool  `json:"enabled,omitempty"`
	APIKey      string `json:"api_key,omitempty"`
}

type saveSMSProviderSettingResponseJSON struct {
	Provider smsProviderSettingJSON `json:"provider"`
}
