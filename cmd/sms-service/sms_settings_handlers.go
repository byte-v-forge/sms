package main

import (
	"errors"
	"net/http"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
)

func (s *dashboardServer) handleSMSSettingsProviderPlugins(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	resp, err := s.smsAdminClient.ListProviderPlugins(r.Context(), &smsinternalv1.ListProviderPluginsRequest{})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}
	if writeProviderError(w, resp.GetError()) {
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (s *dashboardServer) handleSMSSettingsProviders(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		resp, err := s.smsAdminClient.ListProviderConfigs(r.Context(), &smsinternalv1.ListProviderConfigsRequest{IncludeDisabled: true})
		if err != nil {
			writeError(w, http.StatusBadGateway, err)
			return
		}
		if writeProviderError(w, resp.GetError()) {
			return
		}
		writeProtoJSON(w, http.StatusOK, resp)
	case http.MethodPost:
		var req smsinternalv1.UpsertProviderConfigRequest
		if err := readProtoJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		resp, err := s.smsAdminClient.UpsertProviderConfig(r.Context(), &req)
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

func (s *dashboardServer) handleSMSSettingsProvider(w http.ResponseWriter, r *http.Request) {
	id, action, ok := splitSMSPath(r.URL.Path, "/settings/providers/")
	if !ok {
		writeError(w, http.StatusBadRequest, errors.New("provider_key is required"))
		return
	}
	switch {
	case r.Method == http.MethodDelete && action == "":
		s.deleteSMSProviderSetting(w, r, id)
	case r.Method == http.MethodGet && action == "balance":
		s.getSMSProviderBalance(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
