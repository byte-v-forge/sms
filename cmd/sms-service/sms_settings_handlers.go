package main

import (
	"errors"
	"net/http"
)

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
		s.deleteSMSProviderSetting(w, r, id)
	case r.Method == http.MethodGet && action == "balance":
		s.getSMSProviderBalance(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
