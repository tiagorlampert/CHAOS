package utils

import (
	"crypto/rand"
	"regexp"
	"strings"
)

const characters = `0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`

// GenerateRandomString generate a random string based on a given size
func GenerateRandomString(size int) string {
	var bytes = make([]byte, size)
	rand.Read(bytes)
	for i, x := range bytes {
		bytes[i] = characters[x%byte(len(characters))]
	}
	return string(bytes)
}

func NormalizeString(s string) (string, error) {
	re, err := regexp.Compile(`\W`)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(re.ReplaceAllString(s, "")), nil
}
