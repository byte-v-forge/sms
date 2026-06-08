package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	eventbusadapter "github.com/byte-v-forge/sms/internal/adapters/eventbus"
	grpcadapter "github.com/byte-v-forge/sms/internal/adapters/grpc"
	"github.com/byte-v-forge/sms/internal/app"
	smseventcatalog "github.com/byte-v-forge/sms/internal/eventcatalog"
	"github.com/byte-v-forge/sms/internal/platform/eventoutbox"
	"github.com/byte-v-forge/sms/internal/platform/grpcclient"
	"github.com/byte-v-forge/sms/internal/platform/grpchealth"
	"github.com/byte-v-forge/sms/internal/platform/natseventbus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func main() {
	cfg := loadConfig()
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

	httpTimeout := time.Duration(cfg.HTTPTimeoutSeconds) * time.Second
	orderEvents := configuredOrderEvents(platformEventBus != nil && stores.orderOutbox != nil)
	catalogService := app.NewCatalogService(stores.configs, providerRegistry, stores.routeHealth, httpTimeout, cfg.ProviderHTTPProxy, app.SystemClock{})
	orderService := app.NewOrderService(
		stores.orders,
		app.NewConfiguredProviders(providerRegistry, stores.configs, httpTimeout, cfg.ProviderHTTPProxy),
		app.SystemClock{},
		app.RandomIDGenerator{},
		orderEvents,
		hotStream,
		stores.routeHealth,
		stores.codeSecrets,
	)
	adminService := app.NewProviderAdminService(stores.configs, providerRegistry, orderService, stores.orderList, httpTimeout, cfg.ProviderHTTPProxy, hotStream)

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
	startEventWorkers(group, groupCtx, cfg, platformEventBus, stores.orderOutbox, orderService)

	go func() {
		<-groupCtx.Done()
		server.GracefulStop()
	}()

	dashboardConn, err := grpcclient.NewInsecure(grpcclient.SelfTarget(cfg.ListenAddr))
	if err != nil {
		log.Fatalf("connect sms dashboard admin API: %v", err)
	}
	defer dashboardConn.Close()
	errCh := make(chan error, 2)
	startDashboardHTTP(
		groupCtx,
		cfg.DashboardHTTPAddr,
		cfg.DashboardStaticDir,
		smsinternalv1.NewSmsProviderAdminServiceClient(dashboardConn),
		smsv1.NewSmsOrderServiceClient(dashboardConn),
		smsv1.NewSmsCatalogServiceClient(dashboardConn),
		hotStream,
		errCh,
	)

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

func configuredOrderEvents(enabled bool) app.OrderEventSink {
	if !enabled {
		return nil
	}
	return eventbusadapter.NewOrderEventRecorder("sms-service")
}

func startEventWorkers(group *errgroup.Group, ctx context.Context, cfg config, platformEventBus *natseventbus.Bus, orderOutbox *app.PostgresOrderStore, orderService *app.OrderService) {
	if platformEventBus == nil || orderOutbox == nil {
		if strings.TrimSpace(cfg.NATSURL) != "" && orderOutbox == nil {
			log.Print("SMS platform event bus configured without PostgreSQL outbox; using in-process order flow")
		}
		return
	}
	acquireConsumer, err := platformEventBus.PullWorkerForDefinition(cfg.EventStreamName, smseventcatalog.OrderAcquireRequested, 10, 60*time.Second)
	if err != nil {
		log.Fatalf("initialize SMS order acquire worker: %v", err)
	}
	pollConsumer, err := platformEventBus.PullWorkerForDefinition(cfg.EventStreamName, smseventcatalog.OrderPollRequested, 10, 60*time.Second)
	if err != nil {
		log.Fatalf("initialize SMS order poll worker: %v", err)
	}
	cancelConsumer, err := platformEventBus.PullWorkerForDefinition(cfg.EventStreamName, smseventcatalog.OrderCancelRequested, 10, 60*time.Second)
	if err != nil {
		log.Fatalf("initialize SMS order cancel worker: %v", err)
	}
	group.Go(func() error {
		return eventoutbox.RunWorker(ctx, eventoutbox.WorkerConfig{
			Name:      "sms platform event outbox",
			Processor: app.NewOrderOutboxProcessor(orderOutbox, platformEventBus),
			Logf:      log.Printf,
		})
	})
	group.Go(func() error {
		return eventbusadapter.RunOrderAcquireWorker(ctx, acquireConsumer, orderService)
	})
	group.Go(func() error {
		return eventbusadapter.RunOrderPollWorker(ctx, pollConsumer, orderService)
	})
	group.Go(func() error {
		return eventbusadapter.RunOrderCancelWorker(ctx, cancelConsumer, orderService)
	})
}
