package app

import (
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func moneyIsSet(money core.Money) bool {
	return strings.TrimSpace(money.AmountDecimal) != "" || strings.TrimSpace(money.CurrencyCode) != ""
}
