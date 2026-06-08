package protojsonhttp

import (
	"github.com/byte-v-forge/sms/internal/platform/httpx"
	"net/http"
	"strings"

	"github.com/byte-v-forge/sms/internal/platform/protojsonx"

	"google.golang.org/protobuf/proto"
)

func ReadRequest(r *http.Request, dst proto.Message) error {
	defer r.Body.Close()
	raw, err := httpx.ReadLimited(r.Body, httpx.DefaultMaxBodyBytes)
	if err != nil {
		return err
	}
	if len(strings.TrimSpace(string(raw))) == 0 {
		raw = []byte("{}")
	}
	return protojsonx.Unmarshal(raw, dst)
}

func WriteResponse(w http.ResponseWriter, status int, value proto.Message) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	data, err := protojsonx.Marshal(value)
	if err != nil {
		_, _ = w.Write([]byte("{}"))
		return err
	}
	_, writeErr := w.Write(data)
	return writeErr
}
