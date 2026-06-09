package hotstreamnats

import (
	"errors"
	"strings"
)

func normalizeServiceConfig(cfg ServiceConfig) (Config, error) {
	if strings.TrimSpace(cfg.URL) == "" {
		return Config{}, errors.New(requiredServiceConfigMessage(cfg.RequiredMessage))
	}
	service := strings.TrimSpace(cfg.Service)
	clientName := serviceClientName(service, cfg.ClientName)
	return Config{
		URL:        cfg.URL,
		ClientName: clientName,
		Subject:    serviceSubject(service, clientName, cfg.Subject),
		BufferSize: cfg.BufferSize,
	}, nil
}

func normalizeConfig(cfg Config) (Config, error) {
	cfg.URL = strings.TrimSpace(cfg.URL)
	if cfg.URL == "" {
		return Config{}, errors.New("hotstream nats url is required")
	}
	cfg.ClientName = defaultConfigClientName(cfg.ClientName)
	cfg.Subject = defaultConfigSubject(cfg.Subject, cfg.ClientName)
	return cfg, nil
}
