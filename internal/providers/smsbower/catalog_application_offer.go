package smsbower

import "github.com/byte-v-forge/sms/internal/platform/stringx"

func applicationOffersFromList(list []applicationOfferShape) []ApplicationOffer {
	offers := make([]ApplicationOffer, 0, len(list))
	for _, service := range list {
		offers = append(offers, applicationOffer(service.Code, service.Name))
	}
	return offers
}

func applicationOffer(code, name string) ApplicationOffer {
	code = stringx.FirstNonEmpty(code)
	return ApplicationOffer{
		ApplicationKey:     code,
		UpstreamServiceKey: code,
		DisplayName:        stringx.FirstNonEmpty(name, code),
	}
}
