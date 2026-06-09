package hotstreamnats

import (
	"strings"

	"github.com/byte-v-forge/sms/internal/platform/hotstream"
)

func requiredServiceConfigMessage(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "hotstream nats url is required"
	}
	return value
}

func serviceClientName(service string, clientName string) string {
	clientName = strings.TrimSpace(clientName)
	if clientName == "" {
		clientName = service
	}
	if clientName == "" {
		clientName = "byte-v-forge"
	}
	return clientName
}

func serviceSubject(service string, clientName string, subject string) string {
	subject = strings.TrimSpace(subject)
	if subject != "" {
		return subject
	}
	subjectService := service
	if subjectService == "" {
		subjectService = clientName
	}
	return hotstream.ServiceStateSubject(subjectService)
}

func defaultConfigClientName(clientName string) string {
	clientName = strings.TrimSpace(clientName)
	if clientName == "" {
		return "byte-v-forge-hotstream"
	}
	return clientName
}

func defaultConfigSubject(subject string, clientName string) string {
	subject = strings.TrimSpace(subject)
	if subject == "" {
		return hotstream.ServiceStateSubject(clientName)
	}
	return subject
}
