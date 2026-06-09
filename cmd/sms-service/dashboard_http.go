package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/platform/hotstream"
	"golang.org/x/sync/errgroup"
)

type dashboardServer struct {
	smsAdminClient   smsinternalv1.SmsProviderAdminServiceClient
	smsOrderClient   smsv1.SmsOrderServiceClient
	smsCatalogClient smsv1.SmsCatalogServiceClient
	hotstream        hotstream.Subscriber
	staticDir        string
}

func startDashboardHTTP(ctx context.Context, group *errgroup.Group, listenAddr, staticDir string, admin smsinternalv1.SmsProviderAdminServiceClient, orders smsv1.SmsOrderServiceClient, catalog smsv1.SmsCatalogServiceClient, stream hotstream.Subscriber) {
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
	group.Go(func() error {
		log.Printf("sms dashboard BFF listening on %s", listenAddr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("sms dashboard BFF failed: %w", err)
		}
		return nil
	})
}
