package herosms

import (
	"encoding/json"
	"math/big"

	"github.com/byte-v-forge/sms/internal/core"
)

type activationOffersResponse struct {
	Data map[string]map[string]activationOffer `json:"data"`
}

type activationOffer struct {
	Prices struct {
		Default json.RawMessage `json:"default"`
		Retail  json.RawMessage `json:"retail"`
		Min     json.RawMessage `json:"min"`
	} `json:"prices"`
	Counts struct {
		Total        int `json:"total"`
		Physical     int `json:"physical"`
		DefaultPrice int `json:"defaultPrice"`
	} `json:"counts"`
	PriceMap map[string]int `json:"map"`
}

type activationOfferPurchaseTier struct {
	Price          core.Money
	AvailableCount int
	amount         *big.Rat
}
