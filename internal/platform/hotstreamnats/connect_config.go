package hotstreamnats

import (
	"errors"
	"strings"

	"github.com/byte-v-forge/sms/internal/platform/hotstream"
)

func normalizeServiceConfig(cfg ServiceConfig) (Config, error) {
	if strings.TrimSpace(cfg.URL) == "" {
		message := strings.TrimSpace(cfg.RequiredMessage)
		if message == "" {
			message = "hotstream nats url is required"
		}
		return Config{}, errors.New(message)
	}
	clientName := strings.TrimSpace(cfg.ClientName)
	service := strings.TrimSpace(cfg.Service)
	if clientName == "" {
		clientName = service
	}
	if clientName == "" {
		clientName = "byte-v-forge"
	}
	subject := strings.TrimSpace(cfg.Subject)
	if subject == "" {
		subjectService := service
		if subjectService == "" {
			subjectService = clientName
		}
		subject = hotstream.ServiceStateSubject(subjectService)
	}
	return Config{
		URL:        cfg.URL,
		ClientName: clientName,
		Subject:    subject,
		BufferSize: cfg.BufferSize,
	}, nil
}

func normalizeConfig(cfg Config) (Config, error) {
	cfg.URL = strings.TrimSpace(cfg.URL)
	if cfg.URL == "" {
		return Config{}, errors.New("hotstream nats url is required")
	}
	cfg.ClientName = strings.TrimSpace(cfg.ClientName)
	if cfg.ClientName == "" {
		cfg.ClientName = "byte-v-forge-hotstream"
	}
	cfg.Subject = strings.TrimSpace(cfg.Subject)
	if cfg.Subject == "" {
		cfg.Subject = hotstream.ServiceStateSubject(cfg.ClientName)
	}
	return cfg, nil
}
