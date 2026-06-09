package geox

import "strings"

func countryNameOccursInText(text, name string) bool {
	text = " " + text + " "
	name = " " + name + " "
	if strings.Contains(text, name) {
		return true
	}
	return strings.Contains(text, strings.TrimSuffix(name, " ")+"s ")
}
