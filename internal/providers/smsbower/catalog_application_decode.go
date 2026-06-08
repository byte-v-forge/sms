package smsbower

import (
	"encoding/json"

	"github.com/byte-v-forge/sms/internal/core"
)

func decodeApplicationOffers(raw json.RawMessage) ([]ApplicationOffer, error) {
	if len(raw) == 0 {
		return nil, core.NewError(core.CodeUpstreamRejected, "smsbower services list is empty", false)
	}
	if offers, ok := decodeApplicationOfferList(raw); ok {
		return offers, nil
	}
	if offers, ok := decodeApplicationOfferMap(raw); ok {
		return offers, nil
	}
	if offers, ok := decodeApplicationNameMap(raw); ok {
		return offers, nil
	}
	return nil, core.NewError(core.CodeUpstreamRejected, "bad smsbower services list response", false)
}
