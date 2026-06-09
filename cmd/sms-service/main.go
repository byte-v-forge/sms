package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	grpcadapter "github.com/byte-v-forge/sms/internal/adapters/grpc"
	"github.com/byte-v-forge/sms/internal/platform/grpchealth"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("load SMS config: %v", err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	providerRegistry, err := newProviderRegistry()
	if err != nil {
		log.Fatalf("initialize SMS provider registry: %v", err)
	}

	stores, err := newRuntimeStores(ctx, cfg, providerRegistry)
	if err != nil {
		log.Fatalf("initialize SMS stores: %v", err)
	}
	defer stores.close()

	platformEventBus, closePlatformEventBus, err := newPlatformEventBus(ctx, cfg)
	if err != nil {
		log.Fatalf("initialize SMS platform event bus: %v", err)
	}
	defer closePlatformEventBus()
	hotStream, closeHotStream, err := newSMSHotStreamBus(ctx, cfg)
	if err != nil {
		log.Fatalf("initialize SMS hotstream: %v", err)
	}
	defer closeHotStream()

	services := newRuntimeServices(cfg, providerRegistry, stores, platformEventBus, hotStream)

	listener, err := net.Listen("tcp", cfg.ListenAddr)
	if err != nil {
		log.Fatalf("listen %s: %v", cfg.ListenAddr, err)
	}

	server := grpc.NewServer()
	smsv1.RegisterSmsOrderServiceServer(server, grpcadapter.NewOrderServer(services.order))
	smsv1.RegisterSmsCatalogServiceServer(server, grpcadapter.NewCatalogServer(services.catalog))
	smsinternalv1.RegisterSmsProviderAdminServiceServer(server, grpcadapter.NewProviderAdminServer(services.admin))
	grpchealth.RegisterServing(server)
	group, groupCtx := errgroup.WithContext(ctx)
	if err := startEventWorkers(group, groupCtx, cfg, platformEventBus, stores.orderOutbox, services.order); err != nil {
		log.Fatalf("initialize SMS event workers: %v", err)
	}

	go func() {
		<-groupCtx.Done()
		server.GracefulStop()
	}()

	closeDashboard, err := startDashboardBFF(groupCtx, group, cfg, hotStream)
	if err != nil {
		log.Fatalf("initialize sms dashboard BFF: %v", err)
	}
	defer closeDashboard()

	log.Printf("sms-service listening on %s", cfg.ListenAddr)
	group.Go(func() error {
		if err := server.Serve(listener); err != nil && ctx.Err() == nil {
			return err
		}
		return nil
	})
	if err := group.Wait(); err != nil {
		stop()
		log.Fatalf("sms-service failed: %v", err)
	}
}
