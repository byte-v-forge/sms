package app

import (
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func priceCurrencyMismatch(price core.Money, bound core.Money) bool {
	return bound.CurrencyCode != "" && price.CurrencyCode != "" && !strings.EqualFold(bound.CurrencyCode, price.CurrencyCode)
}
