package fivesim

import (
	"context"
	"net/url"

	"github.com/byte-v-forge/sms/internal/core"
)

func (c *Client) SetStatus(ctx context.Context, upstreamOrderID string, action core.ProviderAction) error {
	switch action {
	case core.ActionMarkMessageSent:
		return core.NewError(core.CodeUnsupportedOperation, "5sim does not expose mark-message-sent status", false)
	case core.ActionRequestAdditional:
		return nil
	case core.ActionCompleteOrder:
		var payload order
		return c.getJSON(ctx, "/v1/user/finish/"+url.PathEscape(upstreamOrderID), nil, true, &payload)
	case core.ActionCancelOrder:
		var payload order
		return c.getJSON(ctx, "/v1/user/cancel/"+url.PathEscape(upstreamOrderID), nil, true, &payload)
	default:
		return core.NewError(core.CodeUnsupportedOperation, "unsupported 5sim status action", false)
	}
}
