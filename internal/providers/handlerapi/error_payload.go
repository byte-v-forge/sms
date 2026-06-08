package handlerapi

import "encoding/json"

type handlerAPIErrorPayload struct {
	Title   string                     `json:"title"`
	Detail  string                     `json:"detail"`
	Details string                     `json:"details"`
	Info    map[string]json.RawMessage `json:"info"`
}
