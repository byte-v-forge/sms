package httpsse

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/proto"
)

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
