package main

import (
	"net/http"

	"github.com/byte-v-forge/sms/internal/app"
	"github.com/byte-v-forge/sms/internal/platform/hotstream"
	"github.com/byte-v-forge/sms/internal/platform/httpsse"
)

func (s *dashboardServer) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/settings/provider-plugins", s.handleSMSSettingsProviderPlugins)
	mux.HandleFunc("/settings/providers/", s.handleSMSSettingsProvider)
	mux.HandleFunc("/settings/providers", s.handleSMSSettingsProviders)
	mux.HandleFunc("/price-offers", s.handleSMSPriceOffers)
	mux.HandleFunc("/order-codes", s.handleSMSOrderCodes)
	mux.HandleFunc("/orders/acquire", s.handleSMSOrderAcquire)
	mux.HandleFunc("/orders/", s.handleSMSOrder)
	mux.HandleFunc("/orders", s.handleSMSOrders)
	mux.HandleFunc("/streams/state", s.streamState)
	return mux
}

func (s *dashboardServer) streamState(w http.ResponseWriter, r *http.Request) {
	httpsse.ServeHotStream(w, r, s.hotstream, httpsse.FilterFromRequest(r, hotstream.Filter{
		SourceServices: []string{app.SMSHotStreamSource},
	}), httpsse.ServeOptions{})
}

func (s *dashboardServer) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}
