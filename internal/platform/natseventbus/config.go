package natseventbus

import (
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/platform/eventcatalog"
	"github.com/nats-io/nats.go"
)

const (
	DefaultURL       = nats.DefaultURL
	DefaultStream    = eventcatalog.StreamName
	DefaultSubject   = eventcatalog.StreamSubject
	DefaultFetchWait = 5 * time.Second
)

type Config struct {
	URL        string
	ClientName string
}

func (cfg Config) withDefaults() Config {
	cfg.URL = strings.TrimSpace(cfg.URL)
	if cfg.URL == "" {
		cfg.URL = DefaultURL
	}
	cfg.ClientName = strings.TrimSpace(cfg.ClientName)
	if cfg.ClientName == "" {
		cfg.ClientName = "byte-v-forge"
	}
	return cfg
}
