package app

import (
	"strings"

	"github.com/byte-v-forge/sms/internal/platform/eventoutbox"
)

func nonEmptyRecords(records ...eventoutbox.Record) []eventoutbox.Record {
	out := make([]eventoutbox.Record, 0, len(records))
	for _, record := range records {
		if strings.TrimSpace(record.EventID) == "" {
			continue
		}
		out = append(out, record)
	}
	return out
}
