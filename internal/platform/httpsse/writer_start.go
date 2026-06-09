package httpsse

import "net/http"

func (w *Writer) Start() {
	w.w.Header().Set("Content-Type", "text/event-stream")
	w.w.Header().Set("Cache-Control", "no-cache")
	w.w.Header().Set("Connection", "keep-alive")
	w.w.WriteHeader(http.StatusOK)
	w.Comment("connected")
}
