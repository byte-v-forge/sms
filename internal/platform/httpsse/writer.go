package httpsse

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/protobuf/proto"
)

type Writer struct {
	w       http.ResponseWriter
	flusher http.Flusher
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

func sanitizeLine(value string) string {
	return strings.NewReplacer("\r", " ", "\n", " ").Replace(value)
}
