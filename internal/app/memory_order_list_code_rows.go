package app

import (
	"sort"

	"github.com/byte-v-forge/sms/internal/core"
)

func limitedMemoryOrderCodes(codes []core.OrderCode, limit int) []core.OrderCode {
	codes = append([]core.OrderCode{}, codes...)
	sort.Slice(codes, func(i, j int) bool {
		return codes[i].Code.ReceivedAt.After(codes[j].Code.ReceivedAt)
	})
	if len(codes) > limit {
		codes = codes[:limit]
	}
	return cloneOrderCodes(codes)
}

func cloneOrderCodes(codes []core.OrderCode) []core.OrderCode {
	out := make([]core.OrderCode, 0, len(codes))
	for _, code := range codes {
		out = append(out, core.OrderCode{OrderID: code.OrderID, Code: cloneSMSCode(code.Code)})
	}
	return out
}
