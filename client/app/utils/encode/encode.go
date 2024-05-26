package encode

import (
	"encoding/base64"
	"encoding/json"
)

func Base64Encode(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}

func StringToByte(value string) []byte {
	return []byte(value)
}

func PrettyJson(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func DecodeBase64(value string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
