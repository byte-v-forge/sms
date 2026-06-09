package eventcatalog

import (
	"strings"

	"github.com/byte-v-forge/sms/internal/platform/eventbus"
)

func (definition Definition) ExpectedMessage() eventbus.ExpectedMessage {
	return eventbus.ExpectedMessage{
		Subject:      strings.TrimSpace(definition.Subject),
		EventName:    strings.TrimSpace(definition.EventName),
		EventVersion: strings.TrimSpace(definition.EventVersion),
		PayloadType:  strings.TrimSpace(definition.PayloadType),
	}
}
