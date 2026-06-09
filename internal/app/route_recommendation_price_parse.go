package app

import (
	"math/big"
	"strings"
)

func parseDecimalAmount(value string) (*big.Rat, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, false
	}
	amount, ok := new(big.Rat).SetString(value)
	if !ok {
		return nil, false
	}
	return amount, true
}
