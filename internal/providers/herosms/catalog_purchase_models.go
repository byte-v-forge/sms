package herosms

import "encoding/json"

type purchaseInfoResponse struct {
	Data map[string]purchaseInfo `json:"data"`
}

type purchaseInfo struct {
	Operators            []purchaseOperator `json:"operators"`
	ActivationFinishTime int                `json:"activationFinishTime"`
	UserPrice            json.RawMessage    `json:"userPrice"`
}

type purchaseOperator struct {
	Name             string         `json:"name"`
	LocalName        string         `json:"localName"`
	ActivationsCount int            `json:"activationsCount"`
	FreePriceOffers  map[string]int `json:"freePriceOffers"`
}
