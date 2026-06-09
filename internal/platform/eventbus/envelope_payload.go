package eventbus

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

func envelopePayload(message Message) ([]byte, string, error) {
	if message.Event == nil {
		return nil, "", ErrEmptyEvent
	}
	payload, err := proto.Marshal(message.Event)
	if err != nil {
		return nil, "", fmt.Errorf("marshal event payload: %w", err)
	}
	payloadType := string(message.Event.ProtoReflect().Descriptor().FullName())
	return payload, payloadType, nil
}
