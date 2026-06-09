package httpsse

import (
	"errors"
	"net/http"
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
