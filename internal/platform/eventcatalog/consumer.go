package eventcatalog

import (
	"errors"
)

var ErrEmptyDefinitionConsumerDurable = errors.New("event catalog consumer durable is required")

type ConsumerBinding struct {
	Definition Definition
	Durable    string
}
