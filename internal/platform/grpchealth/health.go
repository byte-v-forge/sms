package grpchealth

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
)

func RegisterServing(registrar grpc.ServiceRegistrar) *health.Server {
	server := health.NewServer()
	healthv1.RegisterHealthServer(registrar, server)
	server.SetServingStatus("", healthv1.HealthCheckResponse_SERVING)
	return server
}

func SetNotServing(server *health.Server) {
	if server == nil {
		return
	}
	server.SetServingStatus("", healthv1.HealthCheckResponse_NOT_SERVING)
}
