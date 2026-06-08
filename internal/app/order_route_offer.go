package app

import (
	"crypto/sha256"
	"encoding/base64"
	"strings"

	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

func publicOfferID(route core.Route) string {
	values := []string{
		normalizeProviderKey(route.ProviderKey),
		strings.TrimSpace(route.ApplicationKey),
		strings.ToUpper(strings.TrimSpace(route.CountryISO2)),
		strings.TrimPrefix(strings.TrimSpace(route.CountryCallingCode), "+"),
		strings.TrimSpace(route.UpstreamServiceKey),
		strings.TrimSpace(route.ProviderCountryID),
		strings.TrimSpace(route.UpstreamProviderID),
	}
	if strings.Join(values, "") == "" {
		return ""
	}
	hash := sha256.Sum256([]byte(strings.Join(values, "\x1f")))
	return "offer_" + base64.RawURLEncoding.EncodeToString(hash[:16])
}

func targetIsZero(target *smsv1.SmsTarget) bool {
	if target == nil {
		return true
	}
	return strings.TrimSpace(target.GetApplicationKey()) == "" &&
		strings.TrimSpace(target.GetCountryIso2()) == "" &&
		strings.TrimSpace(target.GetCountryCallingCode()) == ""
}

func offerRouteRefIsZero(routeRef *smsv1.SmsOfferRouteRef) bool {
	if routeRef == nil {
		return true
	}
	return strings.TrimSpace(routeRef.GetUpstreamServiceKey()) == "" &&
		strings.TrimSpace(routeRef.GetProviderCountryId()) == "" &&
		strings.TrimSpace(routeRef.GetUpstreamProviderId()) == ""
}
