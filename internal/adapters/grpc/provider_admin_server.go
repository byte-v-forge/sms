package grpcadapter

import (
	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/app"
)

type ProviderAdminServer struct {
	smsinternalv1.UnimplementedSmsProviderAdminServiceServer
	service *app.ProviderAdminService
}

func NewProviderAdminServer(service *app.ProviderAdminService) *ProviderAdminServer {
	return &ProviderAdminServer{service: service}
}
