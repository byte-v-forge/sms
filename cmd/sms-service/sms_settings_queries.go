package main

import (
	"errors"
	"net/http"
	"strings"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
)

func (s *dashboardServer) listSMSProviderSettings(r *http.Request) (listSMSProviderSettingsResponseJSON, error) {
	plugins, err := s.smsAdminClient.ListProviderPlugins(r.Context(), &smsinternalv1.ListProviderPluginsRequest{})
	if err != nil {
		return listSMSProviderSettingsResponseJSON{}, err
	}
	if providerErr := plugins.GetError(); providerErr != nil {
		return listSMSProviderSettingsResponseJSON{}, providerError(providerErr)
	}
	configs, err := s.smsAdminClient.ListProviderConfigs(r.Context(), &smsinternalv1.ListProviderConfigsRequest{IncludeDisabled: true})
	if err != nil {
		return listSMSProviderSettingsResponseJSON{}, err
	}
	if providerErr := configs.GetError(); providerErr != nil {
		return listSMSProviderSettingsResponseJSON{}, providerError(providerErr)
	}
	options := make([]smsProviderOptionJSON, 0, len(plugins.GetPlugins()))
	for _, plugin := range plugins.GetPlugins() {
		options = append(options, providerOptionJSON(plugin))
	}
	providers := make([]smsProviderSettingJSON, 0, len(configs.GetConfigs()))
	for _, config := range configs.GetConfigs() {
		providers = append(providers, providerSettingJSON(config))
	}
	return listSMSProviderSettingsResponseJSON{ProviderOptions: options, Providers: providers}, nil
}

func (s *dashboardServer) saveSMSProviderSetting(r *http.Request, req saveSMSProviderSettingRequestJSON) (smsProviderSettingJSON, error) {
	providerKey := strings.TrimSpace(req.ProviderKey)
	if providerKey == "" {
		return smsProviderSettingJSON{}, errors.New("provider_key is required")
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	resp, err := s.smsAdminClient.UpsertProviderConfig(r.Context(), &smsinternalv1.UpsertProviderConfigRequest{
		Config: &smsinternalv1.SmsProviderConfig{
			ProviderKey:      providerKey,
			Enabled:          enabled,
			CredentialSecret: strings.TrimSpace(req.APIKey),
		},
	})
	if err != nil {
		return smsProviderSettingJSON{}, err
	}
	if providerErr := resp.GetError(); providerErr != nil {
		return smsProviderSettingJSON{}, providerError(providerErr)
	}
	return providerSettingJSON(resp.GetConfig()), nil
}
