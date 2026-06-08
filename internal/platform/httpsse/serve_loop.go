package httpsse

import (
	"errors"
	"net/http"
	"time"

	"github.com/byte-v-forge/sms/internal/platform/hotstream"

	observabilityv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/observability/v1"
)

func serveHotStreamLoop(r *http.Request, sse *Writer, sub *hotstream.Subscription, opts ServeOptions) {
	sse.Start()
	sse.Event("", opts.ControlEventName, control(observabilityv1.HotStreamControlKind_HOT_STREAM_CONTROL_KIND_CONNECTED, "connected"))
	ticker := time.NewTicker(opts.Heartbeat)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case event, ok := <-sub.Events:
			if !ok {
				if errors.Is(sub.Err(), hotstream.ErrSlowConsumer) {
					sse.Event("", opts.ControlEventName, control(observabilityv1.HotStreamControlKind_HOT_STREAM_CONTROL_KIND_RESYNC_REQUIRED, "slow consumer; refetch required"))
				}
				return
			}
			sse.Event(event.GetMetadata().GetId(), opts.EventName, event)
		case <-ticker.C:
			sse.Event("", opts.ControlEventName, control(observabilityv1.HotStreamControlKind_HOT_STREAM_CONTROL_KIND_HEARTBEAT, "heartbeat"))
		}
	}
}
