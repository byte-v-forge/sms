package eventcatalog

import (
	"errors"
	"strings"

	"github.com/byte-v-forge/sms/internal/platform/eventbus"
)

var ErrEmptyDefinitionConsumerDurable = errors.New("event catalog consumer durable is required")

type ConsumerBinding struct {
	Definition Definition
	Durable    string
}

func (definition Definition) ExpectedMessage() eventbus.ExpectedMessage {
	return eventbus.ExpectedMessage{
		Subject:      strings.TrimSpace(definition.Subject),
		EventName:    strings.TrimSpace(definition.EventName),
		EventVersion: strings.TrimSpace(definition.EventVersion),
		PayloadType:  strings.TrimSpace(definition.PayloadType),
	}
}

func (definition Definition) DefaultConsumerBinding() ConsumerBinding {
	return ConsumerBinding{Definition: definition, Durable: definition.ConsumerDurable}
}

func (definition Definition) ConsumerBinding(durable string) ConsumerBinding {
	return ConsumerBinding{Definition: definition, Durable: durable}
}

func (binding ConsumerBinding) Subject() string {
	return strings.TrimSpace(binding.Definition.Subject)
}

func (binding ConsumerBinding) DurableName() string {
	durable := strings.TrimSpace(binding.Durable)
	if durable == "" {
		durable = strings.TrimSpace(binding.Definition.ConsumerDurable)
	}
	return durable
}

func (binding ConsumerBinding) Validate() error {
	if binding.Subject() == "" {
		return ErrEmptyDefinitionSubject
	}
	if binding.DurableName() == "" {
		return ErrEmptyDefinitionConsumerDurable
	}
	if strings.TrimSpace(binding.Definition.EventName) == "" {
		return ErrEmptyDefinitionEventName
	}
	if strings.TrimSpace(binding.Definition.EventVersion) == "" {
		return ErrEmptyDefinitionEventVersion
	}
	if strings.TrimSpace(binding.Definition.PayloadType) == "" {
		return ErrEmptyDefinitionPayloadType
	}
	return nil
}
