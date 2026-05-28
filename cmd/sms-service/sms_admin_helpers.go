package main

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
)

func splitSMSPath(path, prefix string) (string, string, bool) {
	tail := strings.Trim(strings.TrimPrefix(path, prefix), "/")
	parts := strings.Split(tail, "/")
	if len(parts) == 0 || strings.TrimSpace(parts[0]) == "" {
		return "", "", false
	}
	id, err := url.PathUnescape(parts[0])
	if err != nil || strings.TrimSpace(id) == "" {
		return "", "", false
	}
	action := ""
	if len(parts) > 1 {
		action = strings.TrimSpace(parts[1])
	}
	return strings.TrimSpace(id), action, true
}

func writeProviderError(w http.ResponseWriter, err *smsinternalv1.ProviderError) bool {
	if err == nil || err.GetPublicError() == nil {
		return false
	}
	writeError(w, http.StatusBadGateway, providerError(err))
	return true
}

func providerError(err *smsinternalv1.ProviderError) error {
	if err == nil {
		return errors.New("sms provider error")
	}
	if err.GetPublicError() == nil {
		return errors.New("sms provider error")
	}
	message := strings.TrimSpace(err.GetPublicError().GetMessage())
	if message == "" {
		message = err.GetPublicError().GetCode().String()
	}
	return errors.New(message)
}
