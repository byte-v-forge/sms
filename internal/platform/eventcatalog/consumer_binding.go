package eventcatalog

import "strings"

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
