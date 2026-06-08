package fivesim

import "encoding/json"

type order struct {
	ID        json.RawMessage `json:"id"`
	CreatedAt string          `json:"created_at"`
	Phone     string          `json:"phone"`
	Operator  string          `json:"operator"`
	Product   string          `json:"product"`
	Price     json.RawMessage `json:"price"`
	Status    string          `json:"status"`
	Expires   string          `json:"expires"`
	SMS       []sms           `json:"sms"`
	Country   string          `json:"country"`
}

type sms struct {
	ID        json.RawMessage `json:"id"`
	CreatedAt string          `json:"created_at"`
	Date      string          `json:"date"`
	Sender    string          `json:"sender"`
	Text      string          `json:"text"`
	Code      string          `json:"code"`
}
