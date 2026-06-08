package httpsse

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/platform/hotstream"

	observabilityv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/observability/v1"
)

func ServeHotStream(w http.ResponseWriter, r *http.Request, subscriber hotstream.Subscriber, filter hotstream.Filter, opts ServeOptions) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if subscriber == nil {
		http.Error(w, "hotstream subscriber is not configured", http.StatusServiceUnavailable)
		return
	}
	sse, err := NewWriter(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sub, err := subscriber.Subscribe(r.Context(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer sub.Close()

	eventName := nonEmpty(opts.EventName, DefaultEventName)
	controlName := nonEmpty(opts.ControlEventName, DefaultControlName)
	heartbeat := opts.Heartbeat
	if heartbeat <= 0 {
		heartbeat = DefaultHeartbeat
	}
	sse.Start()
	sse.Event("", controlName, control(observabilityv1.HotStreamControlKind_HOT_STREAM_CONTROL_KIND_CONNECTED, "connected"))
	ticker := time.NewTicker(heartbeat)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case event, ok := <-sub.Events:
			if !ok {
				if errors.Is(sub.Err(), hotstream.ErrSlowConsumer) {
					sse.Event("", controlName, control(observabilityv1.HotStreamControlKind_HOT_STREAM_CONTROL_KIND_RESYNC_REQUIRED, "slow consumer; refetch required"))
				}
				return
			}
			sse.Event(event.GetMetadata().GetId(), eventName, event)
		case <-ticker.C:
			sse.Event("", controlName, control(observabilityv1.HotStreamControlKind_HOT_STREAM_CONTROL_KIND_HEARTBEAT, "heartbeat"))
		}
	}
}

func nonEmpty(value string, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}
