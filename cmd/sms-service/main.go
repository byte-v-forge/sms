package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/byte-v-forge/common-lib/eventcatalog"
	"github.com/byte-v-forge/common-lib/eventoutbox"
	smsv1 "github.com/byte-v-forge/common-lib/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/common-lib/grpcclient"
	"github.com/byte-v-forge/common-lib/grpchealth"
	"github.com/byte-v-forge/common-lib/redisx"
	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	eventbusadapter "github.com/byte-v-forge/sms/internal/adapters/eventbus"
	grpcadapter "github.com/byte-v-forge/sms/internal/adapters/grpc"
	"github.com/byte-v-forge/sms/internal/app"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func main() {
	cfg := loadConfig()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	configStore, err := app.NewPostgresProviderConfigStore(ctx, cfg.PGDSN)
	if err != nil {
		log.Fatalf("initialize SMS config store: %v", err)
	}
	defer configStore.Close()

	orderHistoryStore, err := app.NewPostgresOrderStore(ctx, cfg.PGDSN)
	if err != nil {
		log.Fatalf("initialize SMS order store: %v", err)
	}
	defer orderHistoryStore.Close()
	redisClient, err := redisx.NewRequiredClient(ctx, cfg.PlatformRedisURL, "PLATFORM_REDIS_URL is required for SMS order state")
	if err != nil {
		log.Fatalf("initialize SMS order redis: %v", err)
	}
	defer redisClient.Close()
	activeStore := app.NewRedisOrderStore(redisx.NewStringStore(redisClient, "sms:order", 30*time.Minute), app.SystemClock{})
	orderStore := app.NewCompositeOrderStore(activeStore, orderHistoryStore)
	routeHealthStore := app.NewRedisRouteHealthStore(redisClient)

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

	httpTimeout := time.Duration(cfg.HTTPTimeoutSeconds) * time.Second
	orderEvents := eventbusadapter.NewOrderEventRecorder("sms-service")
	catalogService := app.NewCatalogService(configStore, routeHealthStore, httpTimeout, cfg.ProviderHTTPProxy, app.SystemClock{})
	orderService := app.NewOrderService(
		orderStore,
		app.NewConfiguredProviders(configStore, httpTimeout, cfg.ProviderHTTPProxy),
		app.SystemClock{},
		app.RandomIDGenerator{},
		orderEvents,
		hotStream,
		routeHealthStore,
	)
	acquireConsumer, err := platformEventBus.PullWorkerConsumer(cfg.EventStreamName, eventcatalog.SMSOrderAcquireRequested.Subject, eventcatalog.SMSOrderAcquireRequested.ConsumerDurable, 10, 60*time.Second)
	if err != nil {
		log.Fatalf("initialize SMS order acquire worker: %v", err)
	}
	pollConsumer, err := platformEventBus.PullWorkerConsumer(cfg.EventStreamName, eventcatalog.SMSOrderPollRequested.Subject, eventcatalog.SMSOrderPollRequested.ConsumerDurable, 10, 60*time.Second)
	if err != nil {
		log.Fatalf("initialize SMS order poll worker: %v", err)
	}
	cancelConsumer, err := platformEventBus.PullWorkerConsumer(cfg.EventStreamName, eventcatalog.SMSOrderCancelRequested.Subject, eventcatalog.SMSOrderCancelRequested.ConsumerDurable, 10, 60*time.Second)
	if err != nil {
		log.Fatalf("initialize SMS order cancel worker: %v", err)
	}
	adminService := app.NewProviderAdminService(configStore, orderService, orderHistoryStore, httpTimeout, cfg.ProviderHTTPProxy, hotStream)

	listener, err := net.Listen("tcp", cfg.ListenAddr)
	if err != nil {
		log.Fatalf("listen %s: %v", cfg.ListenAddr, err)
	}

	server := grpc.NewServer()
	smsv1.RegisterSmsOrderServiceServer(server, grpcadapter.NewOrderServer(orderService))
	smsv1.RegisterSmsCatalogServiceServer(server, grpcadapter.NewCatalogServer(catalogService))
	smsinternalv1.RegisterSmsProviderAdminServiceServer(server, grpcadapter.NewProviderAdminServer(adminService))
	grpchealth.RegisterServing(server)
	group, groupCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		return eventoutbox.RunWorker(groupCtx, eventoutbox.WorkerConfig{
			Name:      "sms platform event outbox",
			Processor: app.NewOrderOutboxProcessor(orderHistoryStore, platformEventBus),
			Logf:      log.Printf,
		})
	})
	group.Go(func() error {
		return eventbusadapter.RunOrderAcquireWorker(groupCtx, acquireConsumer, orderService)
	})
	group.Go(func() error {
		return eventbusadapter.RunOrderPollWorker(groupCtx, pollConsumer, orderService)
	})
	group.Go(func() error {
		return eventbusadapter.RunOrderCancelWorker(groupCtx, cancelConsumer, orderService)
	})

	go func() {
		<-groupCtx.Done()
		server.GracefulStop()
	}()

	dashboardConn, err := grpcclient.NewInsecure(selfTarget(cfg.ListenAddr))
	if err != nil {
		log.Fatalf("connect sms dashboard admin API: %v", err)
	}
	defer dashboardConn.Close()
	errCh := make(chan error, 2)
	startDashboardHTTP(groupCtx, cfg.DashboardHTTPAddr, cfg.DashboardStaticDir, smsinternalv1.NewSmsProviderAdminServiceClient(dashboardConn), hotStream, errCh)

	log.Printf("sms-service listening on %s", cfg.ListenAddr)
	group.Go(func() error {
		if err := server.Serve(listener); err != nil && ctx.Err() == nil {
			return err
		}
		return nil
	})
	group.Go(func() error {
		select {
		case <-groupCtx.Done():
			return nil
		case err := <-errCh:
			return err
		}
	})
	if err := group.Wait(); err != nil {
		stop()
		log.Fatalf("sms-service failed: %v", err)
	}
}
