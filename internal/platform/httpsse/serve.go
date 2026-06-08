package httpsse

import (
	"net/http"

	"github.com/byte-v-forge/sms/internal/platform/hotstream"
)

func ServeHotStream(w http.ResponseWriter, r *http.Request, subscriber hotstream.Subscriber, filter hotstream.Filter, opts ServeOptions) {
	if !validateHotStreamRequest(w, r, subscriber) {
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
	serveHotStreamLoop(r, sse, sub, normalizeServeOptions(opts))
}
