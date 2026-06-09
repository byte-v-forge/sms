package httpsse

import (
	"fmt"
	"strings"
)

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
