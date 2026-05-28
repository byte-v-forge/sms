package grpcadapter

import (
	"context"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/byte-v-forge/sms/internal/app"
	"github.com/byte-v-forge/sms/internal/core"
)

type ProviderAdminServer struct {
	smsinternalv1.UnimplementedSmsProviderAdminServiceServer
	service *app.ProviderAdminService
}

func NewProviderAdminServer(service *app.ProviderAdminService) *ProviderAdminServer {
	return &ProviderAdminServer{service: service}
}

func (s *ProviderAdminServer) ListProviderPlugins(ctx context.Context, request *smsinternalv1.ListProviderPluginsRequest) (*smsinternalv1.ListProviderPluginsResponse, error) {
	plugins, err := s.service.ListProviderPlugins(ctx)
	if err != nil {
		return &smsinternalv1.ListProviderPluginsResponse{Error: toProviderError(err)}, nil
	}
	return &smsinternalv1.ListProviderPluginsResponse{Plugins: plugins}, nil
}

func (s *ProviderAdminServer) UpsertProviderConfig(ctx context.Context, request *smsinternalv1.UpsertProviderConfigRequest) (*smsinternalv1.UpsertProviderConfigResponse, error) {
	config, err := s.service.UpsertProviderConfig(ctx, request.GetConfig())
	if err != nil {
		return &smsinternalv1.UpsertProviderConfigResponse{Error: toProviderError(err)}, nil
	}
	return &smsinternalv1.UpsertProviderConfigResponse{Config: config}, nil
}

func (s *ProviderAdminServer) GetProviderConfig(ctx context.Context, request *smsinternalv1.GetProviderConfigRequest) (*smsinternalv1.GetProviderConfigResponse, error) {
	config, err := s.service.GetProviderConfig(ctx, request.GetProviderKey())
	if err != nil {
		return &smsinternalv1.GetProviderConfigResponse{Error: toProviderError(err)}, nil
	}
	return &smsinternalv1.GetProviderConfigResponse{Config: config}, nil
}

func (s *ProviderAdminServer) ListProviderConfigs(ctx context.Context, request *smsinternalv1.ListProviderConfigsRequest) (*smsinternalv1.ListProviderConfigsResponse, error) {
	configs, err := s.service.ListProviderConfigs(ctx, request.GetIncludeDisabled(), request.GetProviderKey())
	if err != nil {
		return &smsinternalv1.ListProviderConfigsResponse{Error: toProviderError(err)}, nil
	}
	return &smsinternalv1.ListProviderConfigsResponse{Configs: configs}, nil
}

func (s *ProviderAdminServer) DeleteProviderConfig(ctx context.Context, request *smsinternalv1.DeleteProviderConfigRequest) (*smsinternalv1.DeleteProviderConfigResponse, error) {
	if err := s.service.DeleteProviderConfig(ctx, request.GetProviderKey()); err != nil {
		return &smsinternalv1.DeleteProviderConfigResponse{Error: toProviderError(err)}, nil
	}
	return &smsinternalv1.DeleteProviderConfigResponse{}, nil
}

func (s *ProviderAdminServer) GetProviderBalance(ctx context.Context, request *smsinternalv1.GetProviderBalanceRequest) (*smsinternalv1.GetProviderBalanceResponse, error) {
	balance, err := s.service.GetProviderBalance(ctx, request.GetProviderKey())
	if err != nil {
		return &smsinternalv1.GetProviderBalanceResponse{Error: toProviderError(err)}, nil
	}
	return &smsinternalv1.GetProviderBalanceResponse{Balance: toProtoMoney(balance)}, nil
}

func (s *ProviderAdminServer) ListOrders(ctx context.Context, request *smsinternalv1.ListOrdersRequest) (*smsinternalv1.ListOrdersResponse, error) {
	orders, err := s.service.ListOrders(ctx, request.GetIncludeFinal(), int(request.GetLimit()))
	if err != nil {
		return &smsinternalv1.ListOrdersResponse{Error: toProviderError(err)}, nil
	}
	out := make([]*smsinternalv1.SmsOrderView, 0, len(orders))
	for _, order := range orders {
		out = append(out, toOrderView(order))
	}
	return &smsinternalv1.ListOrdersResponse{Orders: out}, nil
}

func (s *ProviderAdminServer) CancelOrder(ctx context.Context, request *smsinternalv1.CancelProviderOrderRequest) (*smsinternalv1.CancelProviderOrderResponse, error) {
	order, err := s.service.CancelOrder(ctx, request.GetOrderId(), request.GetRequestId())
	if err != nil {
		return &smsinternalv1.CancelProviderOrderResponse{Order: toOrderView(order), Error: toProviderError(err)}, nil
	}
	return &smsinternalv1.CancelProviderOrderResponse{Order: toOrderView(order)}, nil
}

func toOrderView(order core.Order) *smsinternalv1.SmsOrderView {
	if order.ID == "" {
		return nil
	}
	return &smsinternalv1.SmsOrderView{
		Order:       toProtoOrder(order),
		ProviderKey: order.ProviderKey,
	}
}

func toProviderError(err error) *smsinternalv1.ProviderError {
	if err == nil {
		return nil
	}
	return &smsinternalv1.ProviderError{PublicError: toProtoError(err)}
}
