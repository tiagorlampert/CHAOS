package network

import (
	"encoding/base64"
)

func DecodeToBytes(encData string) []byte {
	decData, _ := base64.URLEncoding.DecodeString(encData)
	return decData
}
