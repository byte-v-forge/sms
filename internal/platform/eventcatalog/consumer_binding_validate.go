package eventcatalog

import "strings"

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
