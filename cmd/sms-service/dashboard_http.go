package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/app"
	"github.com/byte-v-forge/sms/internal/platform/hotstream"
	"github.com/byte-v-forge/sms/internal/platform/httpsse"
	"github.com/byte-v-forge/sms/internal/platform/protojsonhttp"
	"google.golang.org/protobuf/proto"
)

type dashboardServer struct {
	smsAdminClient   smsinternalv1.SmsProviderAdminServiceClient
	smsOrderClient   smsv1.SmsOrderServiceClient
	smsCatalogClient smsv1.SmsCatalogServiceClient
	hotstream        hotstream.Subscriber
	staticDir        string
}

func startDashboardHTTP(ctx context.Context, listenAddr, staticDir string, admin smsinternalv1.SmsProviderAdminServiceClient, orders smsv1.SmsOrderServiceClient, catalog smsv1.SmsCatalogServiceClient, stream hotstream.Subscriber, errCh chan<- error) {
	if strings.TrimSpace(listenAddr) == "" {
		return
	}
	if strings.TrimSpace(staticDir) == "" {
		staticDir = "/app/dashboard/sms"
	}
	dashboard := &dashboardServer{smsAdminClient: admin, smsOrderClient: orders, smsCatalogClient: catalog, hotstream: stream, staticDir: staticDir}
	mux := http.NewServeMux()
	mux.Handle("/api/sms/", http.StripPrefix("/api/sms", dashboard.routes()))
	mux.HandleFunc("/healthz", dashboard.handleHealth)
	mux.Handle("/", noCacheSPAFileServer(staticDir))
	server := &http.Server{Addr: listenAddr, Handler: withCORS(mux), ReadHeaderTimeout: 5 * time.Second}
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()
	go func() {
		log.Printf("sms dashboard BFF listening on %s", listenAddr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("sms dashboard BFF failed: %w", err)
		}
	}()
}

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

func readProtoJSON(r *http.Request, dst proto.Message) error {
	return protojsonhttp.ReadRequest(r, dst)
}

func writeProtoJSON(w http.ResponseWriter, status int, value proto.Message) {
	_ = protojsonhttp.WriteResponse(w, status, value)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func requireHTTPMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method == method {
		return true
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
	return false
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PATCH,PUT,DELETE,OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func noCacheSPAFileServer(dir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		path := filepath.Join(dir, filepath.Clean(r.URL.Path))
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			http.ServeFile(w, r, path)
			return
		}
		http.ServeFile(w, r, filepath.Join(dir, "index.html"))
	})
}
