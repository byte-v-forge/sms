package eventbus

import "strings"

func Attributes(pairs ...string) map[string]string {
	attrs := map[string]string{}
	for i := 0; i+1 < len(pairs); i += 2 {
		attrs = WithAttribute(attrs, pairs[i], pairs[i+1])
	}
	if len(attrs) == 0 {
		return nil
	}
	return attrs
}

func WithAttribute(attrs map[string]string, key string, value string) map[string]string {
	key = strings.TrimSpace(key)
	if key == "" {
		return attrs
	}
	if attrs == nil {
		attrs = map[string]string{}
	}
	attrs[key] = strings.TrimSpace(value)
	return attrs
}

func WithNonEmptyAttribute(attrs map[string]string, key string, value string) map[string]string {
	if strings.TrimSpace(value) == "" {
		return attrs
	}
	return WithAttribute(attrs, key, value)
}
