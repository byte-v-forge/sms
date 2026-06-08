package grpcadapter

import (
	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/app"
	"github.com/byte-v-forge/sms/internal/core"
)

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
		OfferRef:                app.PublicOfferRefFromRoute(offer.Route),
	}
}
