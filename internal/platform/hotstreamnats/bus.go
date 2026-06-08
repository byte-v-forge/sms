package hotstreamnats

import (
	"github.com/byte-v-forge/sms/internal/platform/hotstream"
	"github.com/nats-io/nats.go"
)

type Config struct {
	URL        string
	ClientName string
	Subject    string
	BufferSize int
}

type ServiceConfig struct {
	URL             string
	Service         string
	ClientName      string
	Subject         string
	BufferSize      int
	RequiredMessage string
}

type Bus struct {
	conn    *nats.Conn
	hub     *hotstream.Hub
	subject string
	nodeID  string
	sub     *nats.Subscription
}
