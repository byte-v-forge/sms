package smsbower

import "encoding/json"

type getNumberV2Payload struct {
	OrderID          json.RawMessage `json:"activationId"`
	PhoneNumber      json.RawMessage `json:"phoneNumber"`
	OrderCost        json.RawMessage `json:"activationCost"`
	CanGetAnotherSMS json.RawMessage `json:"canGetAnotherSms"`
	OrderTime        json.RawMessage `json:"activationTime"`
}
