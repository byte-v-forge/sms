package httpsse

import (
	"net/http"

	"github.com/byte-v-forge/sms/internal/platform/hotstream"
)

func validateHotStreamRequest(w http.ResponseWriter, r *http.Request, subscriber hotstream.Subscriber) bool {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return false
	}
	if subscriber == nil {
		http.Error(w, "hotstream subscriber is not configured", http.StatusServiceUnavailable)
		return false
	}
	return true
}
