package geox

import "strings"

func alphaTokens(value string) []string {
	fields := strings.FieldsFunc(value, func(r rune) bool {
		return !((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z'))
	})
	out := []string{}
	for _, field := range fields {
		field = strings.ToUpper(strings.TrimSpace(field))
		if len(field) == 2 || len(field) == 3 {
			out = append(out, field)
		}
	}
	return out
}
