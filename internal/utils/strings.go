package utils

import (
	"regexp"
	"strings"
)

func NormalizeString(s string) (string, error) {
	re, err := regexp.Compile(`\W`)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(re.ReplaceAllString(s, "")), nil
}
