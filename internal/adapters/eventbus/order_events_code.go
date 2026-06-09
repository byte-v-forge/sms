package eventbusadapter

import (
	"context"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/app"
	"github.com/byte-v-forge/sms/internal/core"
	"github.com/byte-v-forge/sms/internal/platform/eventbus"
	"github.com/byte-v-forge/sms/internal/platform/eventcatalog"
	"github.com/byte-v-forge/sms/internal/platform/eventoutbox"
)

func (b *OrderEventRecorder) recordCodeReceived(ctx context.Context, order core.Order, code core.SMSCode) (eventoutbox.Record, error) {
	metadata := b.metadata(
		eventcatalog.SMSCodeReceived.EventName,
		eventcatalog.SMSCodeReceived.Subject,
		eventbus.StableEventID("code-received-", order.ID, eventTimeSuffix(code.ReceivedAt)),
		order.ID,
		code.ReceivedAt,
	)
	return b.record(ctx, eventcatalog.SMSCodeReceived, &smsv1.SmsCodeReceivedEvent{
		Metadata: metadata,
		OrderId:  order.ID,
		Code:     app.PublicCode(&code),
	}, metadata, orderAttributes(order))
}
