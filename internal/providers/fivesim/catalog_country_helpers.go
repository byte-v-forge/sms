package fivesim

func firstMapKey(values map[string]int) string {
	for key := range values {
		return key
	}
	return ""
}

func trimPlus(value string) string {
	if len(value) > 0 && value[0] == '+' {
		return value[1:]
	}
	return value
}
