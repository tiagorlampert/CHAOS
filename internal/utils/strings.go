package utils

import (
	"regexp"
)

func SanitizeString(s string) string {
	re := regexp.MustCompile(`\W`)
	return string(re.ReplaceAll([]byte(s), []byte("")))
}

func SanitizeUrl(original string) string {
	urlRegExp := regexp.MustCompile(`[^a-zA-Z0-9-_/:.,?&@=#%]`)
	return string(urlRegExp.ReplaceAll([]byte(original), []byte("")))
}
