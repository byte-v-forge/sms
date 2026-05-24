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
	config, err := s.service.GetProviderConfig(ctx, request.GetProviderConfigId())
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
	if err := s.service.DeleteProviderConfig(ctx, request.GetProviderConfigId()); err != nil {
		return &smsinternalv1.DeleteProviderConfigResponse{Error: toProviderError(err)}, nil
	}
	return &smsinternalv1.DeleteProviderConfigResponse{}, nil
}

func (s *ProviderAdminServer) ListRouteOptions(ctx context.Context, request *smsinternalv1.ListRouteOptionsRequest) (*smsinternalv1.ListRouteOptionsResponse, error) {
	options, err := s.service.ListRouteOptions(ctx, request.GetProviderConfigId(), request.GetProviderKey())
	if err != nil {
		return &smsinternalv1.ListRouteOptionsResponse{Error: toProviderError(err)}, nil
	}
	return &smsinternalv1.ListRouteOptionsResponse{Options: options}, nil
}

func (s *ProviderAdminServer) UpsertRouteProfile(ctx context.Context, request *smsinternalv1.UpsertRouteProfileRequest) (*smsinternalv1.UpsertRouteProfileResponse, error) {
	profile, err := s.service.UpsertRouteProfile(ctx, request.GetProfile())
	if err != nil {
		return &smsinternalv1.UpsertRouteProfileResponse{Error: toProviderError(err)}, nil
	}
	return &smsinternalv1.UpsertRouteProfileResponse{Profile: profile}, nil
}

func (s *ProviderAdminServer) GetRouteProfile(ctx context.Context, request *smsinternalv1.GetRouteProfileRequest) (*smsinternalv1.GetRouteProfileResponse, error) {
	profile, err := s.service.GetRouteProfile(ctx, request.GetProfileKey())
	if err != nil {
		return &smsinternalv1.GetRouteProfileResponse{Error: toProviderError(err)}, nil
	}
	return &smsinternalv1.GetRouteProfileResponse{Profile: profile}, nil
}

func (s *ProviderAdminServer) ListRouteProfiles(ctx context.Context, request *smsinternalv1.ListRouteProfilesRequest) (*smsinternalv1.ListRouteProfilesResponse, error) {
	profiles, err := s.service.ListRouteProfiles(ctx, request.GetIncludeDisabled())
	if err != nil {
		return &smsinternalv1.ListRouteProfilesResponse{Error: toProviderError(err)}, nil
	}
	return &smsinternalv1.ListRouteProfilesResponse{Profiles: profiles}, nil
}

func (s *ProviderAdminServer) DeleteRouteProfile(ctx context.Context, request *smsinternalv1.DeleteRouteProfileRequest) (*smsinternalv1.DeleteRouteProfileResponse, error) {
	if err := s.service.DeleteRouteProfile(ctx, request.GetProfileKey()); err != nil {
		return &smsinternalv1.DeleteRouteProfileResponse{Error: toProviderError(err)}, nil
	}
	return &smsinternalv1.DeleteRouteProfileResponse{}, nil
}

func (s *ProviderAdminServer) GetProviderBalance(ctx context.Context, request *smsinternalv1.GetProviderBalanceRequest) (*smsinternalv1.GetProviderBalanceResponse, error) {
	balance, err := s.service.GetProviderBalance(ctx, request.GetProviderConfigId())
	if err != nil {
		return &smsinternalv1.GetProviderBalanceResponse{Error: toProviderError(err)}, nil
	}
	return &smsinternalv1.GetProviderBalanceResponse{Balance: toProtoMoney(balance)}, nil
}

func (s *ProviderAdminServer) ListActivations(ctx context.Context, request *smsinternalv1.ListActivationsRequest) (*smsinternalv1.ListActivationsResponse, error) {
	activations, err := s.service.ListActivations(ctx, request.GetIncludeFinal(), int(request.GetLimit()))
	if err != nil {
		return &smsinternalv1.ListActivationsResponse{Error: toProviderError(err)}, nil
	}
	views := make([]*smsinternalv1.SmsActivationView, 0, len(activations))
	for _, activation := range activations {
		views = append(views, toActivationView(activation))
	}
	return &smsinternalv1.ListActivationsResponse{Activations: views}, nil
}

func (s *ProviderAdminServer) CancelActivation(ctx context.Context, request *smsinternalv1.CancelProviderActivationRequest) (*smsinternalv1.CancelProviderActivationResponse, error) {
	activation, err := s.service.CancelActivation(ctx, request.GetActivationId(), request.GetRequestId())
	if err != nil {
		return &smsinternalv1.CancelProviderActivationResponse{Activation: toActivationView(activation), Error: toProviderError(err)}, nil
	}
	return &smsinternalv1.CancelProviderActivationResponse{Activation: toActivationView(activation)}, nil
}

func toActivationView(activation core.Activation) *smsinternalv1.SmsActivationView {
	if activation.ID == "" {
		return nil
	}
	return &smsinternalv1.SmsActivationView{
		Activation:           toProtoActivation(activation),
		LatestCode:           toProtoCode(activation.Code),
		ProviderConfigId:     activation.ProviderConfigID,
		ProviderKey:          activation.ProviderKey,
		UpstreamActivationId: activation.UpstreamActivationID,
		UpstreamOperator:     activation.UpstreamOperator,
	}
}

func toProviderError(err error) *smsinternalv1.ProviderError {
	if err == nil {
		return nil
	}
	return &smsinternalv1.ProviderError{PublicError: toProtoError(err)}
}
