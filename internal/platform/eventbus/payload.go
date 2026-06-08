package eventbus

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

func UnmarshalPayload(message ReceivedMessage, event proto.Message) error {
	if event == nil {
		return ErrEmptyEvent
	}
	if message.Envelope == nil || len(message.Envelope.GetPayload()) == 0 {
		return ErrEmptyPayload
	}
	if err := proto.Unmarshal(message.Envelope.GetPayload(), event); err != nil {
		return fmt.Errorf("unmarshal event payload: %w", err)
	}
	return nil
}
