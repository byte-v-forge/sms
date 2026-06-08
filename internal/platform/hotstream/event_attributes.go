package hotstream

func CleanAttributes(input map[string]string) map[string]string {
	if len(input) == 0 {
		return nil
	}
	out := map[string]string{}
	for key, value := range input {
		key = cleanEventText(key)
		value = cleanEventText(value)
		if key == "" || value == "" {
			continue
		}
		out[key] = value
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
