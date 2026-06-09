package main

import (
	"context"
	"strings"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/platform/grpcclient"
	"github.com/byte-v-forge/sms/internal/platform/hotstream"
	"golang.org/x/sync/errgroup"
)

func startDashboardBFF(ctx context.Context, group *errgroup.Group, cfg config, stream hotstream.Subscriber) (func(), error) {
	if strings.TrimSpace(cfg.DashboardHTTPAddr) == "" {
		return noopClose, nil
	}
	conn, err := grpcclient.NewInsecure(grpcclient.SelfTarget(cfg.ListenAddr))
	if err != nil {
		return nil, err
	}
	startDashboardHTTP(
		ctx,
		group,
		cfg.DashboardHTTPAddr,
		cfg.DashboardStaticDir,
		smsinternalv1.NewSmsProviderAdminServiceClient(conn),
		smsv1.NewSmsOrderServiceClient(conn),
		smsv1.NewSmsCatalogServiceClient(conn),
		stream,
	)
	return func() { _ = conn.Close() }, nil
}
