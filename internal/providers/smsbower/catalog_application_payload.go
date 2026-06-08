package smsbower

import (
	"context"
	"encoding/json"
)

type applicationsPayload struct {
	Status   string          `json:"status"`
	Services json.RawMessage `json:"services"`
}

func (c *Client) getApplicationsPayload(ctx context.Context) (applicationsPayload, error) {
	result, err := c.api.Do(ctx, "getServicesList", nil)
	if err != nil {
		return applicationsPayload{}, err
	}
	var payload applicationsPayload
	if err := decodeJSONObject(result, &payload); err != nil {
		return applicationsPayload{}, err
	}
	return payload, nil
}
