package main

import (
	"errors"
	"net/http"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/app"
)

func (s *dashboardServer) handleSMSOrders(w http.ResponseWriter, r *http.Request) {
	if !requireHTTPMethod(w, r, http.MethodGet) {
		return
	}
	resp, err := s.smsAdminClient.ListOrders(r.Context(), smsListOrdersRequest(r))
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}
	if writeProviderError(w, resp.GetError()) {
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (s *dashboardServer) handleSMSPriceOffers(w http.ResponseWriter, r *http.Request) {
	if !requireHTTPMethod(w, r, http.MethodGet) {
		return
	}
	if s.smsCatalogClient == nil {
		writeError(w, http.StatusServiceUnavailable, errors.New("sms catalog service is not configured"))
		return
	}
	resp, err := s.smsCatalogClient.ListSmsPriceOffers(r.Context(), smsPriceOffersRequest(r))
	if err != nil {
		writeProtoJSON(w, http.StatusOK, &smsv1.ListSmsPriceOffersResponse{Error: app.PublicError(err)})
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (s *dashboardServer) handleSMSApplications(w http.ResponseWriter, r *http.Request) {
	if !requireHTTPMethod(w, r, http.MethodGet) {
		return
	}
	if s.smsCatalogClient == nil {
		writeError(w, http.StatusServiceUnavailable, errors.New("sms catalog service is not configured"))
		return
	}
	resp, err := s.smsCatalogClient.ListSmsApplications(r.Context(), smsApplicationsRequest(r))
	if err != nil {
		writeProtoJSON(w, http.StatusOK, &smsv1.ListSmsApplicationsResponse{Error: app.PublicError(err)})
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (s *dashboardServer) handleSMSCountries(w http.ResponseWriter, r *http.Request) {
	if !requireHTTPMethod(w, r, http.MethodGet) {
		return
	}
	if s.smsCatalogClient == nil {
		writeError(w, http.StatusServiceUnavailable, errors.New("sms catalog service is not configured"))
		return
	}
	resp, err := s.smsCatalogClient.ListSmsCountries(r.Context(), smsCountriesRequest(r))
	if err != nil {
		writeProtoJSON(w, http.StatusOK, &smsv1.ListSmsCountriesResponse{Error: app.PublicError(err)})
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (s *dashboardServer) handleSMSOrderCodes(w http.ResponseWriter, r *http.Request) {
	if !requireHTTPMethod(w, r, http.MethodGet) {
		return
	}
	resp, err := s.smsAdminClient.ListOrderCodes(r.Context(), smsListOrderCodesRequest(r))
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}
	if writeProviderError(w, resp.GetError()) {
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}
