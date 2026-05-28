package grpcadapter

import (
	"context"
	"sort"

	smsv1 "github.com/byte-v-forge/common-lib/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/app"
	"github.com/byte-v-forge/sms/internal/core"
)

type CatalogServer struct {
	smsv1.UnimplementedSmsCatalogServiceServer
	service *app.CatalogService
}

func NewCatalogServer(service *app.CatalogService) *CatalogServer {
	return &CatalogServer{service: service}
}

func (s *CatalogServer) ListSmsProviders(ctx context.Context, _ *smsv1.ListSmsProvidersRequest) (*smsv1.ListSmsProvidersResponse, error) {
	providers, err := s.service.ListProviders(ctx)
	if err != nil {
		return &smsv1.ListSmsProvidersResponse{Error: toProtoError(err)}, nil
	}
	out := make([]*smsv1.SmsProviderInfo, 0, len(providers))
	for _, provider := range providers {
		capabilities := provider.GetCapabilities()
		out = append(out, &smsv1.SmsProviderInfo{
			ProviderKey:             provider.GetProviderKey(),
			DisplayName:             provider.GetDisplayName(),
			SupportsBalance:         capabilities.GetSupportsBalance(),
			SupportsCatalog:         capabilities.GetSupportsCatalog(),
			SupportsAdditionalCode:  capabilities.GetSupportsAdditionalCode(),
			RequiresMarkMessageSent: capabilities.GetRequiresMarkMessageSent(),
		})
	}
	return &smsv1.ListSmsProvidersResponse{Providers: out}, nil
}

func (s *CatalogServer) ListSmsPriceOffers(ctx context.Context, request *smsv1.ListSmsPriceOffersRequest) (*smsv1.ListSmsPriceOffersResponse, error) {
	offers, err := s.service.ListPriceOffers(ctx, core.RouteOfferQuery{
		ApplicationKey:     request.GetApplicationKey(),
		CountryISO2:        request.GetCountryIso2(),
		CountryCallingCode: request.GetCountryCallingCode(),
		ProviderKey:        request.GetProviderKey(),
	})
	if err != nil {
		return &smsv1.ListSmsPriceOffersResponse{Error: toProtoError(err)}, nil
	}
	return &smsv1.ListSmsPriceOffersResponse{Offers: toProtoPriceOffers(offers)}, nil
}

func (s *CatalogServer) ListSmsApplications(ctx context.Context, request *smsv1.ListSmsApplicationsRequest) (*smsv1.ListSmsApplicationsResponse, error) {
	offers, err := s.service.ListPriceOffers(ctx, core.RouteOfferQuery{ProviderKey: request.GetProviderKey()})
	if err != nil {
		return &smsv1.ListSmsApplicationsResponse{Error: toProtoError(err)}, nil
	}
	items := map[string]*smsv1.SmsApplicationInfo{}
	for _, offer := range offers {
		if offer.ApplicationKey == "" {
			continue
		}
		items[offer.ApplicationKey] = &smsv1.SmsApplicationInfo{ApplicationKey: offer.ApplicationKey, DisplayName: offer.ApplicationName}
	}
	out := make([]*smsv1.SmsApplicationInfo, 0, len(items))
	for _, item := range items {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].GetApplicationKey() < out[j].GetApplicationKey() })
	return &smsv1.ListSmsApplicationsResponse{Applications: out}, nil
}

func (s *CatalogServer) ListSmsCountries(ctx context.Context, request *smsv1.ListSmsCountriesRequest) (*smsv1.ListSmsCountriesResponse, error) {
	offers, err := s.service.ListPriceOffers(ctx, core.RouteOfferQuery{ProviderKey: request.GetProviderKey()})
	if err != nil {
		return &smsv1.ListSmsCountriesResponse{Error: toProtoError(err)}, nil
	}
	items := map[string]*smsv1.SmsCountry{}
	for _, offer := range offers {
		key := offer.CountryISO2
		if key == "" {
			key = offer.CountryCallingCode
		}
		if key == "" {
			continue
		}
		items[key] = &smsv1.SmsCountry{CountryIso2: offer.CountryISO2, Name: offer.CountryName, CountryCallingCode: offer.CountryCallingCode}
	}
	out := make([]*smsv1.SmsCountry, 0, len(items))
	for _, item := range items {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].GetCountryIso2()+out[i].GetCountryCallingCode() < out[j].GetCountryIso2()+out[j].GetCountryCallingCode()
	})
	return &smsv1.ListSmsCountriesResponse{Countries: out}, nil
}

func toProtoPriceOffers(offers []core.RouteOffer) []*smsv1.SmsPriceOffer {
	out := make([]*smsv1.SmsPriceOffer, 0, len(offers))
	for _, offer := range offers {
		out = append(out, toProtoPriceOffer(offer))
	}
	return out
}

func toProtoPriceOffer(offer core.RouteOffer) *smsv1.SmsPriceOffer {
	return &smsv1.SmsPriceOffer{
		ApplicationKey:          offer.ApplicationKey,
		ApplicationName:         offer.ApplicationName,
		CountryIso2:             offer.CountryISO2,
		CountryName:             offer.CountryName,
		CountryCallingCode:      offer.CountryCallingCode,
		Price:                   toProtoMoney(offer.Price),
		AvailableCount:          int32(offer.AvailableCount),
		ProviderKey:             offer.ProviderKey,
		ProviderDisplayName:     offer.ProviderDisplayName,
		SupportsCancel:          offer.SupportsCancel,
		SupportsAdditionalCode:  offer.SupportsAdditionalCode,
		RequiresMarkMessageSent: offer.RequiresMarkMessageSent,
		ObservedAt:              toProtoTime(offer.ObservedAt),
		UpstreamProviderId:      offer.UpstreamProviderID,
		UpstreamProviderName:    offer.UpstreamProviderName,
		AcquireParams:           app.PublicAcquireParamsFromRoute(offer.Route),
	}
}
