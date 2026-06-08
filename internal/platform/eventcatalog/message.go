package eventcatalog

import (
	"strings"

	commonv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/common/v1"
	"github.com/byte-v-forge/sms/internal/platform/eventbus"
	"google.golang.org/protobuf/proto"
)

func (definition Definition) NewMessage(
	event proto.Message,
	metadata *commonv1.EventMetadata,
	attributes map[string]string,
) (eventbus.Message, error) {
	if err := definition.ValidateEvent(event); err != nil {
		return eventbus.Message{}, err
	}
	if err := definition.ValidateMetadata(metadata); err != nil {
		return eventbus.Message{}, err
	}
	return eventbus.Message{
		Subject:    strings.TrimSpace(definition.Subject),
		Event:      event,
		Metadata:   metadata,
		Extensions: attributes,
	}, nil
}
