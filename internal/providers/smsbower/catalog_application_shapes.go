package smsbower

import (
	"encoding/json"

	"github.com/byte-v-forge/sms/internal/platform/stringx"
)

type applicationOfferShape struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func decodeApplicationOfferList(raw json.RawMessage) ([]ApplicationOffer, bool) {
	var list []applicationOfferShape
	if err := json.Unmarshal(raw, &list); err != nil {
		return nil, false
	}
	return applicationOffersFromList(list), true
}

func decodeApplicationOfferMap(raw json.RawMessage) ([]ApplicationOffer, bool) {
	var byCode map[string]applicationOfferShape
	if err := json.Unmarshal(raw, &byCode); err != nil {
		return nil, false
	}
	offers := make([]ApplicationOffer, 0, len(byCode))
	for code, item := range byCode {
		offers = append(offers, applicationOffer(stringx.FirstNonEmpty(item.Code, code), stringx.FirstNonEmpty(item.Name, code)))
	}
	return offers, true
}

func decodeApplicationNameMap(raw json.RawMessage) ([]ApplicationOffer, bool) {
	var names map[string]string
	if err := json.Unmarshal(raw, &names); err != nil {
		return nil, false
	}
	offers := make([]ApplicationOffer, 0, len(names))
	for code, name := range names {
		offers = append(offers, applicationOffer(code, name))
	}
	return offers, true
}
