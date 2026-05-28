package main

import (
	"errors"
	"net/http"
	"strings"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
)

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

func (s *dashboardServer) handleSMSSettingsProviders(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		settings, err := s.listSMSProviderSettings(r)
		if err != nil {
			writeError(w, http.StatusBadGateway, err)
			return
		}
		writeJSON(w, http.StatusOK, settings)
	case http.MethodPost:
		var req saveSMSProviderSettingRequestJSON
		if err := readJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		provider, err := s.saveSMSProviderSetting(r, req)
		if err != nil {
			writeError(w, http.StatusBadGateway, err)
			return
		}
		writeJSON(w, http.StatusOK, saveSMSProviderSettingResponseJSON{Provider: provider})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *dashboardServer) handleSMSSettingsProvider(w http.ResponseWriter, r *http.Request) {
	id, action, ok := splitSMSPath(r.URL.Path, "/settings/providers/")
	if !ok {
		writeError(w, http.StatusBadRequest, errors.New("provider_key is required"))
		return
	}
	switch {
	case r.Method == http.MethodDelete && action == "":
		resp, err := s.smsAdminClient.DeleteProviderConfig(r.Context(), &smsinternalv1.DeleteProviderConfigRequest{ProviderKey: id})
		if err != nil {
			writeError(w, http.StatusBadGateway, err)
			return
		}
		if writeProviderError(w, resp.GetError()) {
			return
		}
		writeProtoJSON(w, http.StatusOK, resp)
	case r.Method == http.MethodGet && action == "balance":
		resp, err := s.smsAdminClient.GetProviderBalance(r.Context(), &smsinternalv1.GetProviderBalanceRequest{ProviderKey: id})
		if err != nil {
			writeError(w, http.StatusBadGateway, err)
			return
		}
		if writeProviderError(w, resp.GetError()) {
			return
		}
		writeProtoJSON(w, http.StatusOK, resp)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

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
			ProviderKey:       providerKey,
			Enabled:           enabled,
			CredentialSecret:  strings.TrimSpace(req.APIKey),
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
