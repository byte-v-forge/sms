package eventbus

import (
	"context"
	"strings"
)

func applyHandlerResult(ctx context.Context, message ReceivedMessage, result HandlerResult, logf LogFunc) {
	label := strings.TrimSpace(result.Label)
	switch result.Action {
	case MessageActionNak:
		NakMessageDelay(ctx, message, result.Delay, defaultResultLabel(label, "nak event"), logf)
	case MessageActionTerm:
		TermMessage(ctx, message, defaultResultLabel(label, "terminate event"), logf)
	default:
		AckMessage(ctx, message, defaultResultLabel(label, "ack event"), logf)
	}
}

func defaultResultLabel(label string, fallback string) string {
	if label == "" {
		return fallback
	}
	return label
}
