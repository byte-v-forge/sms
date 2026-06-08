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
		routeText(route.ApplicationKey),
		routeCountryISO2(route.CountryISO2),
		routeCallingCode(route.CountryCallingCode),
		routeText(route.UpstreamServiceKey),
		routeText(route.ProviderCountryID),
		routeText(route.UpstreamProviderID),
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
	return routeText(target.GetApplicationKey()) == "" &&
		routeCountryISO2(target.GetCountryIso2()) == "" &&
		routeCallingCode(target.GetCountryCallingCode()) == ""
}

func offerRouteRefIsZero(routeRef *smsv1.SmsOfferRouteRef) bool {
	if routeRef == nil {
		return true
	}
	return routeText(routeRef.GetUpstreamServiceKey()) == "" &&
		routeText(routeRef.GetProviderCountryId()) == "" &&
		routeText(routeRef.GetUpstreamProviderId()) == ""
}
