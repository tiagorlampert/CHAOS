package utils

func ParseCheckboxBoolean(value string) bool {
	if value == "true" {
		return true
	}
	return false
}
