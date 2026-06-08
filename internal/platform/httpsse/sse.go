package httpsse

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/platform/hotstream"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	observabilityv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/observability/v1"
)

const (
	DefaultEventName   = "hotstream"
	DefaultControlName = "hotstream.control"
	DefaultHeartbeat   = 15 * time.Second
)

type Writer struct {
	w       http.ResponseWriter
	flusher http.Flusher
}

type ServeOptions struct {
	EventName        string
	ControlEventName string
	Heartbeat        time.Duration
	Logf             func(string, ...any)
}

func NewWriter(w http.ResponseWriter) (*Writer, error) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, errors.New("streaming is not supported")
	}
	return &Writer{w: w, flusher: flusher}, nil
}

func (w *Writer) Start() {
	w.w.Header().Set("Content-Type", "text/event-stream")
	w.w.Header().Set("Cache-Control", "no-cache")
	w.w.Header().Set("Connection", "keep-alive")
	w.w.WriteHeader(http.StatusOK)
	w.Comment("connected")
}

func (w *Writer) Event(id string, name string, message proto.Message) {
	id = strings.TrimSpace(id)
	name = strings.TrimSpace(name)
	if id != "" {
		_, _ = fmt.Fprintf(w.w, "id: %s\n", sanitizeLine(id))
	}
	if name != "" {
		_, _ = fmt.Fprintf(w.w, "event: %s\n", sanitizeLine(name))
	}
	_, _ = fmt.Fprintf(w.w, "data: %s\n\n", protoJSON(message))
	w.flusher.Flush()
}

func (w *Writer) Comment(text string) {
	text = strings.TrimSpace(text)
	if text == "" {
		text = "keepalive"
	}
	_, _ = fmt.Fprintf(w.w, ": %s\n\n", sanitizeLine(text))
	w.flusher.Flush()
}

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

func control(kind observabilityv1.HotStreamControlKind, message string) *observabilityv1.HotStreamControlEvent {
	return &observabilityv1.HotStreamControlEvent{Kind: kind, Message: message, OccurredAt: timestamppb.Now()}
}

func protoJSON(message proto.Message) string {
	data, err := (protojson.MarshalOptions{UseProtoNames: true}).Marshal(message)
	if err != nil {
		fallback, _ := (protojson.MarshalOptions{UseProtoNames: true}).Marshal(control(observabilityv1.HotStreamControlKind_HOT_STREAM_CONTROL_KIND_ERROR, err.Error()))
		return string(fallback)
	}
	return string(data)
}

func nonEmpty(value string, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}

func sanitizeLine(value string) string {
	return strings.NewReplacer("\r", " ", "\n", " ").Replace(value)
}
