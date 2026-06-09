package hotstreamnats

import (
	"strings"

	"github.com/byte-v-forge/sms/internal/platform/hotstream"
)

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
