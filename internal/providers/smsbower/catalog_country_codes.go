package smsbower

import (
	"strings"

	"github.com/byte-v-forge/sms/internal/platform/geox"
)

func smsbowerCountryCodes(name string) (string, string) {
	normalized := strings.TrimSpace(name)
	if normalized == "" {
		return "", ""
	}
	iso2 := geox.CountryAlpha2InText(normalized)
	if iso2 == "" {
		return "", ""
	}
	return iso2, geox.CountryCallingCode(iso2)
}
