package hotstreamnats

import (
	"strings"

	"github.com/byte-v-forge/sms/internal/platform/hotstream"
)

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
