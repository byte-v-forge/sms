package main

import (
	"net/http"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
)

func (s *dashboardServer) deleteSMSProviderSetting(w http.ResponseWriter, r *http.Request, providerKey string) {
	resp, err := s.smsAdminClient.DeleteProviderConfig(r.Context(), &smsinternalv1.DeleteProviderConfigRequest{ProviderKey: providerKey})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}
	if writeProviderError(w, resp.GetError()) {
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (s *dashboardServer) getSMSProviderBalance(w http.ResponseWriter, r *http.Request, providerKey string) {
	resp, err := s.smsAdminClient.GetProviderBalance(r.Context(), &smsinternalv1.GetProviderBalanceRequest{ProviderKey: providerKey})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}
	if writeProviderError(w, resp.GetError()) {
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}
