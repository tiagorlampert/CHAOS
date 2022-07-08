package utils

import "encoding/base64"

// EncodeBase64 returns a encoded string to base64
func EncodeBase64(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}

// DecodeBase64 returns a decoded string from base64
func DecodeBase64(value string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
