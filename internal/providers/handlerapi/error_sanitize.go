package handlerapi

import "github.com/byte-v-forge/sms/internal/providers/providerhttp"

func truncateHandlerAPIErrorMessage(message string) string {
	return providerhttp.SanitizeErrorText(message)
}
