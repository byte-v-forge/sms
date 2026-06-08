package hotstream

import "strings"

func ServiceStateSubject(service string) string {
	service = strings.Trim(strings.ToLower(strings.TrimSpace(service)), ".")
	if service == "" {
		service = "platform"
	}
	return SubjectPrefix + "." + service + ".state"
}
