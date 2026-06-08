package grpcadapter

import (
	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/app"
)

type OrderServer struct {
	smsv1.UnimplementedSmsOrderServiceServer
	service *app.OrderService
}

func NewOrderServer(service *app.OrderService) *OrderServer {
	return &OrderServer{service: service}
}
